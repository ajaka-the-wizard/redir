CREATE OR REPLACE FUNCTION pseudo_encrypt(VALUE int) RETURNS int AS $$
DECLARE 
    l1 int;
    l2 int;
    r1 int;
    r2 int;
    i int:=0;
BEGIN
    l1:=(VALUE >> 16) & 65535;
    r1:= VALUE & 65535;
    WHILE i < 3 LOOP
        l2 := r1;
        r2 := l1;
        l1 := l2;
        r1 := r2;
        i := i + 1;
    END LOOP;
    RETURN ((l1 << 16) | r1);
END;
$$ LANGUAGE plpgsql STRICT IMMUTABLE;