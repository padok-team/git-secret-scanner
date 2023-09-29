import tempfile


# directory name of the temporary directory used by the tool
TEMP_DIR_NAME = 'github.padok.git-secret-scanner'

# default path to use for cloning when None is specified
DEFAULT_CLONE_PATH = f'{tempfile.gettempdir()}/{TEMP_DIR_NAME}'
