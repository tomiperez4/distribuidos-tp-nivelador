import fcntl
from contextlib import contextmanager

@contextmanager
def _file_lock(path, mode):
    with open(path, "a+") as lock_file:
        fcntl.flock(lock_file, mode)
        try:
            yield lock_file
        finally:
            fcntl.flock(lock_file, fcntl.LOCK_UN)