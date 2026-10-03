import zlib, struct, sys

def ppm_to_png(path_in, path_out):
    ppm = open(path_in, "rb").read()
    parts = ppm.split(b"\n", 3)
    w, h = map(int, parts[1].split())
    maxv = int(parts[2])
    pix = parts[3]
    # PNG with IDAT zlib of raw scanlines (filter 0)
    rows = b""
    for row in range(h):
        rows += b"\x00" + pix[row*w*3:(row+1)*w*3]
    comp = zlib.compress(rows, 6)
    def chunk(tag, data):
        c = struct.pack(">I", len(data)) + tag + data
        c += struct.pack(">I", zlib.crc32(tag + data) & 0xffffffff)
        return c
    png = b"\x89PNG\r\n\x1a\n"
    png += chunk(b"IHDR", struct.pack(">IIBBBBB", w, h, 8, 2, 0, 0, 0))
    png += chunk(b"IDAT", comp)
    png += chunk(b"IEND", b"")
    open(path_out, "wb").write(png)
    print("png written:", path_out, os.path.getsize(path_out) if (os := __import__("os")) else "(size n/a)")

ppm_to_png(r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\120.ppm",
           r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\120.png")
