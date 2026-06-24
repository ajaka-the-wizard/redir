package redir

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	defaultMaxRetries = 3
	defaultTimeout    = 30 * time.Second
)

var (
	ErrIncompleteBatch = errors.New("redir: upload batch incomplete")
	ErrNoFiles         = errors.New("redir: no files provided")
)

type Config struct {
	BaseURL    string
	ProductID  int
	APIKey     string
	Auto       *bool
	MaxRetries int
	HTTPClient *http.Client
}

type Client struct {
	baseURL    string
	productID  int
	apiKey     string
	auto       bool
	maxRetries int
	httpClient *http.Client
}

type File struct {
	Name        string
	ContentType string
	Body        io.ReadSeeker
}

type Media struct {
	PublicKey string    `json:"public_key"`
	BatchID   uuid.UUID `json:"batch_id"`
	SeqID     int       `json:"seq_id"`
	UserID    uuid.UUID `json:"user_id"`
	FileSize  int64     `json:"file_size"`
	Status    string    `json:"status"`
	FileType  string    `json:"file_type"`
	Active    bool      `json:"active"`
	Public    bool      `json:"public"`
	FileName  string    `json:"file_name"`
	MimeType  string    `json:"mime_type"`
	Hits      int       `json:"hits"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UploadOptions struct {
	BatchID *uuid.UUID
	Auto    *bool
}

type UploadResult struct {
	Success       bool
	Committed     bool
	BatchID       uuid.UUID
	Media         []Media
	MissingSeqIDs []int
	Attempts      int
}

type IncompleteBatchError struct {
	BatchID       uuid.UUID
	MissingSeqIDs []int
	Attempts      int
}

func (e *IncompleteBatchError) Error() string {
	return fmt.Sprintf("redir: batch %s incomplete after %d attempt(s), missing seq ids: %v", e.BatchID, e.Attempts, e.MissingSeqIDs)
}

func (e *IncompleteBatchError) Is(target error) bool {
	return target == ErrIncompleteBatch
}

func Bool(v bool) *bool {
	return &v
}

func New(cfg Config) (*Client, error) {
	if strings.TrimSpace(cfg.BaseURL) == "" {
		return nil, errors.New("redir: base url is required")
	}
	if cfg.ProductID <= 0 {
		return nil, errors.New("redir: product id is required")
	}
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, errors.New("redir: api key is required")
	}

	auto := true
	if cfg.Auto != nil {
		auto = *cfg.Auto
	}
	maxRetries := cfg.MaxRetries
	if maxRetries <= 0 {
		maxRetries = defaultMaxRetries
	}
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}

	return &Client{
		baseURL:    strings.TrimRight(cfg.BaseURL, "/"),
		productID:  cfg.ProductID,
		apiKey:     cfg.APIKey,
		auto:       auto,
		maxRetries: maxRetries,
		httpClient: httpClient,
	}, nil
}

func FileFromPath(path string) (File, error) {
	f, err := os.Open(path)
	if err != nil {
		return File{}, err
	}
	return File{
		Name: filepath.Base(path),
		Body: f,
	}, nil
}

func FileFromBytes(name, contentType string, data []byte) File {
	return File{
		Name:        name,
		ContentType: contentType,
		Body:        bytes.NewReader(data),
	}
}

func (f File) Close() error {
	closer, ok := f.Body.(io.Closer)
	if !ok {
		return nil
	}
	return closer.Close()
}

func (c *Client) Ping(ctx context.Context) error {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/client/ping", nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeAPIError(resp)
}

func (c *Client) Upload(ctx context.Context, files []File, opts *UploadOptions) (*UploadResult, error) {
	if len(files) == 0 {
		return nil, ErrNoFiles
	}
	for i := range files {
		if strings.TrimSpace(files[i].Name) == "" {
			return nil, fmt.Errorf("redir: file at index %d has no name", i)
		}
		if files[i].Body == nil {
			return nil, fmt.Errorf("redir: file %q has no body", files[i].Name)
		}
	}

	batchID := uuid.New()
	if opts != nil && opts.BatchID != nil {
		batchID = *opts.BatchID
	}
	auto := c.auto
	if opts != nil && opts.Auto != nil {
		auto = *opts.Auto
	}

	wanted := make([]int, len(files))
	pending := make([]int, len(files))
	for i := range files {
		wanted[i] = i
		pending[i] = i
	}

	mediaBySeq := map[int]Media{}
	attemptsAllowed := 1
	if auto {
		attemptsAllowed += c.maxRetries
	}

	result := &UploadResult{BatchID: batchID}
	for attempt := 1; attempt <= attemptsAllowed && len(pending) > 0; attempt++ {
		media, err := c.uploadOnce(ctx, batchID, files, pending)
		result.Attempts = attempt
		if err != nil {
			return result, err
		}
		for _, m := range media {
			mediaBySeq[m.SeqID] = m
		}
		pending = missingSeqIDs(wanted, mediaBySeq)
	}

	result.Media = orderedMedia(wanted, mediaBySeq)
	result.MissingSeqIDs = pending
	if len(pending) > 0 {
		return result, &IncompleteBatchError{
			BatchID:       batchID,
			MissingSeqIDs: append([]int(nil), pending...),
			Attempts:      result.Attempts,
		}
	}

	if err := c.Commit(ctx, batchID); err != nil {
		return result, err
	}
	result.Success = true
	result.Committed = true
	return result, nil
}

func (c *Client) Commit(ctx context.Context, batchID uuid.UUID) error {
	req, err := c.newRequest(ctx, http.MethodPut, "/api/v1/client/commit/"+batchID.String(), nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeAPIError(resp)
}

func (c *Client) uploadOnce(ctx context.Context, batchID uuid.UUID, files []File, seqIDs []int) ([]Media, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for _, seqID := range seqIDs {
		file := files[seqID]
		if _, err := file.Body.Seek(0, io.SeekStart); err != nil {
			return nil, fmt.Errorf("redir: could not rewind file %q: %w", file.Name, err)
		}
		contentType := file.ContentType
		if contentType == "" {
			contentType = detectContentType(file.Body)
			if _, err := file.Body.Seek(0, io.SeekStart); err != nil {
				return nil, fmt.Errorf("redir: could not rewind file %q: %w", file.Name, err)
			}
		}
		part, err := writer.CreatePart(filePartHeader(file.Name, contentType, seqID))
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(part, file.Body); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/v1/client/upload", &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Batch-ID", batchID.String())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := decodeAPIError(resp); err != nil {
		return nil, err
	}

	var data uploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Media, nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Product", strconv.Itoa(c.productID))
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	return req, nil
}

type uploadResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Media   []Media `json:"media"`
}

type apiErrorResponse struct {
	Message string   `json:"message"`
	Error   string   `json:"error"`
	Errors  []string `json:"errors"`
}

func decodeAPIError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	var data apiErrorResponse
	_ = json.NewDecoder(resp.Body).Decode(&data)
	msg := data.Message
	if msg == "" {
		msg = data.Error
	}
	if msg == "" && len(data.Errors) > 0 {
		msg = strings.Join(data.Errors, ", ")
	}
	if msg == "" {
		msg = http.StatusText(resp.StatusCode)
	}
	return fmt.Errorf("redir: api returned %d: %s", resp.StatusCode, msg)
}

func filePartHeader(name, contentType string, seqID int) textproto.MIMEHeader {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, escapeQuotes(name)))
	h.Set("Content-Type", contentType)
	h.Set("X-Sequential-ID", strconv.Itoa(seqID))
	return h
}

func detectContentType(r io.Reader) string {
	var sample [512]byte
	n, _ := r.Read(sample[:])
	if n == 0 {
		return "application/octet-stream"
	}
	return http.DetectContentType(sample[:n])
}

func escapeQuotes(s string) string {
	return strings.NewReplacer("\\", "\\\\", `"`, "\\\"").Replace(s)
}

func missingSeqIDs(wanted []int, mediaBySeq map[int]Media) []int {
	missing := make([]int, 0)
	for _, seqID := range wanted {
		if _, ok := mediaBySeq[seqID]; !ok {
			missing = append(missing, seqID)
		}
	}
	return missing
}

func orderedMedia(wanted []int, mediaBySeq map[int]Media) []Media {
	media := make([]Media, 0, len(mediaBySeq))
	for _, seqID := range wanted {
		if m, ok := mediaBySeq[seqID]; ok {
			media = append(media, m)
		}
	}
	sort.SliceStable(media, func(i, j int) bool {
		return media[i].SeqID < media[j].SeqID
	})
	return media
}
