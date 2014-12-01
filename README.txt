Copyright (c) 2014, Waitman Gobble <ns@waitman.net>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer. 
2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


This software is a simple example command line music player which uses JSON output 
from the modified metaflac program.

*This software depends on a modified metaflac from https://github.com/waitman/flac*

First,
Edit flacplayer.go and adjust location of metaflac, and mplayer

Before running the software the first time you need to create a postgresql 
database. If you have postgresql server running on localhost, and your 
username has access to create databases and tables, then the following
command should create the table

# createdb goflac

Change the db connection string if you need to access a remote server 
and/or adjust SSL or authentication values.

1) Build

# go build

2) Run

./goflac

The first time goflac runs, it will create two tables. 
tracks and images. If you specify a path as argument 
it will scan that directory and subdirectories for *.flac
files and import the meta information into the database.


Example Run

% ./goflac /xj/waitman/Music/
/xj/waitman/Music/01 - Canto de Ossanha.flac already in database.
/xj/waitman/Music/02 - Inside My Head.flac already in database.
/xj/waitman/Music/03 - Rio.flac already in database.
/xj/waitman/Music/04 - Reza.flac already in database.
/xj/waitman/Music/05 - Canto de Iemanja.flac already in database.
/xj/waitman/Music/06 - Juanito.flac already in database.
/xj/waitman/Music/07 - Casa forte.flac already in database.
/xj/waitman/Music/Thriller/01-Wanna Be Startin' Somethin'.flac already in database.
/xj/waitman/Music/Thriller/02-Baby Be Mine.flac already in database.
/xj/waitman/Music/Thriller/03-The Girl Is Mine (feat. Paul McCartney).flac already in database.
/xj/waitman/Music/Thriller/04-Thriller.flac already in database.
/xj/waitman/Music/Thriller/05-Beat It.flac already in database.
/xj/waitman/Music/Thriller/06-Billie Jean.flac already in database.
/xj/waitman/Music/Thriller/07-Human Nature.flac already in database.
/xj/waitman/Music/Thriller/08-P.Y.T. (Pretty Young Thing).flac already in database.
/xj/waitman/Music/Thriller/09-The Lady In My Life.flac already in database.

#       Chan    Rate    BPS     Path
---     ----    ------  ---     ----
16      2       176400  24      /xj/waitman/Music/Thriller/09-The Lady In My Life.flac
15      2       176400  24      /xj/waitman/Music/Thriller/08-P.Y.T. (Pretty Young Thing).flac
14      2       176400  24      /xj/waitman/Music/Thriller/07-Human Nature.flac
13      2       176400  24      /xj/waitman/Music/Thriller/06-Billie Jean.flac
12      2       176400  24      /xj/waitman/Music/Thriller/05-Beat It.flac
11      2       176400  24      /xj/waitman/Music/Thriller/04-Thriller.flac
10      2       176400  24      /xj/waitman/Music/Thriller/03-The Girl Is Mine (feat. Paul McCartney).flac
9       2       176400  24      /xj/waitman/Music/Thriller/02-Baby Be Mine.flac
8       2       176400  24      /xj/waitman/Music/Thriller/01-Wanna Be Startin' Somethin'.flac
7       2       88200   24      /xj/waitman/Music/07 - Casa forte.flac
6       2       88200   24      /xj/waitman/Music/06 - Juanito.flac
5       2       88200   24      /xj/waitman/Music/05 - Canto de Iemanja.flac
4       2       88200   24      /xj/waitman/Music/04 - Reza.flac
3       2       88200   24      /xj/waitman/Music/03 - Rio.flac
2       2       88200   24      /xj/waitman/Music/02 - Inside My Head.flac
1       2       88200   24      /xj/waitman/Music/01 - Canto de Ossanha.flac

Enter the Track Number: (Control-C to exit) 
3

Now Playing: /xj/waitman/Music/03 - Rio.flac

TITLE                                           (Rio) - Rio
ALBUM                                           Jazzanova
ISRC                                            DEPI81400575
GENRE                                           Contemporary Jazz
YEAR                                            2014
DATE                                            2811
PART_OF_A_SET                                   1/1
ITUNESCOMPILATION                               0
COMPOSER                                        Menescal, Roberto
COMPOSER                                        Freire,
BAND                                            Jazzanova
INVOLVED_PEOPLE_LIST                            ComposerMenescal, Roberto
INVOLVED_PEOPLE_LIST                            ComposerFreire
INVOLVED_PEOPLE_LIST                            EnsembleJazzanova
TAGGING_TIME                                    2014-11-19T10:33:39+01:00
COPYRIGHT                                       eClassical
PUBLISHER                                       MPS Records
PUBLISHER                                       eClassical
COMMERCIAL                                      http://www.eclassical.com
COMMENT                                         Downloaded from eClassical.com. From album 4250644876967
TRACKNUMBER                                     3/7

 (Press Return to Quit)


Enter the Track Number: (Control-C to exit) 


