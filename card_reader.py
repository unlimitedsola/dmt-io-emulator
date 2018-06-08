import serial
import os
import sys

card_uid = b"\xf2\x8a\x0d\x00"
card_id = "1UZV VIXZ OJFL 2U12 FEDY".replace(' ', '')
card_id_part_1 = bytearray(card_id[:16])
card_id_part_2 = bytearray(card_id[16:20]) + b"\x00" * 12

responses = {
    # is there any card? -> NOPE!
    b"\x02\x00\x02\x31\x30\x03\x02": b"\x02\x00\x03\x31\x30\x4e\x03\x4d",
    # eject card -> OK!
    b"\x02\x00\x02\x32\x30\x03\x01": b"\x02\x00\x03",
    # find card -> We got one!
    b"\x02\x00\x02\x35\x30\x03\x06": b"\x02\x00\x03\x35\x30\x59\x03\x5e",
    # get UID -> return our UID!
    b"\x02\x00\x02\x35\x31\x03\x07": b"\x02\x00\x07\x35\x31\x59" + card_uid + b"\x03\x2e",
    # auth key A at sector 0 (key = 37 21 53 6a 72 40) -> You got it!
    b"\x02\x00\x09\x35\x32\x00\x37\x21\x53\x6a\x72\x40\x03\x12": b"\x02\x00\x04\x35\x32\x00\x59\x03\x5b",
    # read sector 0 block 1 -> Our card ID
    b"\x02\x00\x04\x35\x33\x00\x01\x03\x02": b"\x02\x00\x15\x35\x33\x00\x01\x59" + card_id_part_1 + b"\x03\x54",
    b"\x02\x00\x04\x35\x33\x00\x02\x03\x01": b"\x02\x00\x15\x35\x33\x00\x02\x59" + card_id_part_2 + b"\x03\x57"
}

if __name__ == '__main__':
    com = serial.Serial("COM9")
    acc = ""
    while (True):
        in_code = com.read()
        acc += in_code
        print(acc)
        if in_code == b'\x03':
            bcc = com.read()
            acc += bcc
            print('acc=', acc)
            if acc in responses:
                com.write(b'\x06')
                if com.read() == b'\x05':
                    com.write(responses.get(acc))
                    print('write=', responses.get(acc))
            else:
                com.write(b'\x15')
            acc = ""
