from lottery import Bet

BET_HEADER_SIZE = 14
DATE_BYTE_SIZE = 4
MAX_BYTE_SIZE = 255

def decode_bet(payload: bytes, agency_id: int) -> tuple[Bet, int]:
    if len(payload) < BET_HEADER_SIZE:
        raise Exception("decode bet: header too short")

    firstname_len = payload[12]
    lastname_len = payload[13]

    total_len = BET_HEADER_SIZE + firstname_len + lastname_len

    if len(payload) < total_len:
        raise Exception("decode bet: payload too short")

    document = int.from_bytes(payload[:4], byteorder='big')
    birthdate_int = int.from_bytes(payload[4:8], byteorder='big')
    birthdate = __unpack_date__(birthdate_int)
    number = int.from_bytes(payload[8:12], byteorder='big')


    offset = 14
    firstname = payload[offset:offset + firstname_len].decode('utf-8')
    offset += firstname_len
    lastname = payload[offset:offset + lastname_len].decode('utf-8')

    return Bet(agency_id, firstname, lastname, document, birthdate, number), total_len

def encode_bet(bet: Bet) -> bytes:
    firstname = bet.first_name.encode('utf-8')
    lastname = bet.last_name.encode('utf-8')

    if len(firstname) > MAX_BYTE_SIZE or len(lastname) > MAX_BYTE_SIZE:
        raise Exception("encode bet: first/last name too long")

    firstname_len = len(firstname).to_bytes(1, byteorder='big')
    lastname_len = len(lastname).to_bytes(1, byteorder='big')

    document = bet.document.to_bytes(4, byteorder='big')
    birthdate = __pack_date__(bet.birthdate)
    number = bet.number.to_bytes(4, byteorder='big')


    return b"".join([
        document,
        birthdate,
        number,
        firstname_len,
        lastname_len,
        firstname,
        lastname
    ])


def __unpack_date__(date: int) -> str:
    year = date // 10000
    month = (date // 100) % 100
    day = date % 100

    return f"{year}-{month}-{day}"

def __pack_date__(date: str) -> bytes:
    date = date.split('-')
    year = int(date[0])
    month = int(date[1])
    day = int(date[2])

    birthdate = year * 10000 + month * 100 + day
    return birthdate.to_bytes(4, byteorder='big')