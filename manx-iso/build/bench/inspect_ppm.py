import struct, sys
ppm = open(r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\120.ppm","rb").read()
# header: P6\n<w h>\n<max>\n
parts = ppm.split(b"\n", 3)
w, h = map(int, parts[1].split())
maxv = int(parts[2])
pix = parts[3]
print(f"w={w} h={h} max={maxv} pixel-bytes={len(pix)}")
# bright/dark histogram coarse: count near-black and near-white
dark = light = 0
for row in range(0, h, 8):
    base = row * w * 3
    for col in range(0, w, 8):
        i = base + col*3
        r,g,b = pix[i], pix[i+1], pix[i+2]
        if r < 40 and g < 40 and b < 40:
            dark += 1
        elif r > 220 and g > 220 and b > 220:
            light += 1
print(f"sampled dark={dark} light={light}")
