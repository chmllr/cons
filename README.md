cons
====

    Usage: cons --dir <PATH> [OPTIONS] <COMMAND>
    
    cons helps keeping track of file changes in a directory.  If no directory
    parameter was provided, the current directory is used.
    
    Avaliable commands:
    
    seal:
    	Seals all existing files with their md5 hashes into the directory index.
    	It does not make any mutating operations on the directory!
    
    verify (accepts option --deep):
    	Checks existing file structure against the index. This command can detect 
    	missing, modified, duplicated and new files. If option 'deep' is provided, 
    	compares files based on the MD5 hash.
