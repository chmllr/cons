cons
====

    Usage: cons --dir <PATH> [OPTIONS] <COMMAND>
    
    cons keeps track of file changes in a directory. If no directory parameter was
    provided, the current directory is used. cons creates a file index inside the
    `.index.csv' file in CSV format. If certain files should be ignored, create
    a file name `.consignore' inside the tracked directory with one regular
    expression per line, matching the file names to be ignored.
    
    Avaliable commands:
    
    seal:
    	Seals all existing files with their MD5 hashes into the directory index.
    	It does not make any mutating operations on the directory!
    
    verify (accepts option --deep):
    	Checks existing file structure against the index. This command can detect 
    	missing, modified, duplicated and new files. If option 'deep' is provided, 
    	compares files based on the MD5 hash.
