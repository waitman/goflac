package main

import (
  "path/filepath"
  "os"
  "flag"
  "fmt"
  "strings"
_ "github.com/lib/pq"
  "database/sql"
  "log"
  "encoding/json"
  "os/exec"
  "bufio"
  "strconv"
)

type Image struct {
	pictureType int
	mimeType string
	description string
	width int
	height int
	depth int
	dataLength int
	imageData string
}

type Samples struct {
	PointNumber int `json:"Point Number"`
	SampleNumber int `json:"Sample Number"`
	StreamOffset int `json:"Stream Offset"`
	FrameSamples int `json:"Frame Samples"`
}

type Block struct {
	BlockNumber int `json:"Block Number"`
	BlockType   int `json:"Block Type"`
	BlockLength int `json:"Block Length"`
	IsLast bool `json:"Is Last"`
	MinBlocksize int `json:"Min Blocksize"`
	MaxBlocksize int `json:"Max Blocksize"`
	MinFramesize int `json:"Min Framesize"`
	MaxFramesize int `json:"Max Framesize"`
	SampleRate int `json:"Sample Rate"`
	Channels int `json:"Channels"`
	BitsPerSample int `json:"Bits Per Sample"`
	TotalSamples int `json:"Total Samples"`
	MD5Signature string `json:"MD5 Signature"`
	SeekPoints int `json:"Seek Points"`
	SeekRecords []Samples `json:"Seek Data"`
	VendorString string `json:"Vendor String"`
	NumberOfComments int `json:"Number of Comments"`
	Comments []string `json:"Comments"`
	PictureType int `json:"Picture Type"`
	MIMEType string `json:"MIME Type"`
	Description string `json:"Description"`
	Width int `json:"Width"`
	Height int `json:"Height"`
	Depth int `json:"Depth"`
	DataLength int `json:"DataLength"`
	ImageData string `json:"Image Data"`	
}

type Rec struct {
	BlockID int   `json:"Block ID"`
	Detail  Block `json:"Block"`
}

func visit(path string, f os.FileInfo, err error) error {
	if strings.Contains(path,"flac") != false {
		metaflacargs := []string{"-c","/xj/waitman/flac/src/metaflac/metaflac --output-json --list \"" + path + "\""};
		cmd := exec.Command("/bin/sh",metaflacargs...);
		output, err := cmd.CombinedOutput()

		if err!=nil {
			fmt.Printf(fmt.Sprint(err))
			log.Fatal("oh no")
		}

		var recs []Rec
		err = json.Unmarshal([]byte(output), &recs)
		if err !=nil {
			log.Fatal(err)
		}
	
		/* BlockType 0 */
		var md5Signature string
		var minBlockSize int
		var maxBlockSize int
		var minFrameSize int
		var maxFrameSize int
		var sampleRate int
		var channels int
		var bitsPerSample int
		var totalSamples int

		/* BlockType 4*/
		var vendorString string
		//var numberOfComments int
		var comments string

		/* BlockType 6*/

		var images []Image

		for i := range recs {
			switch (recs[i].Detail.BlockType) {
				case 0: 
					md5Signature = recs[i].Detail.MD5Signature
					minBlockSize = recs[i].Detail.MinBlocksize
					maxBlockSize = recs[i].Detail.MaxBlocksize
					minFrameSize = recs[i].Detail.MinFramesize
					maxFrameSize = recs[i].Detail.MaxFramesize
					sampleRate   = recs[i].Detail.SampleRate
					channels     = recs[i].Detail.Channels
					bitsPerSample= recs[i].Detail.BitsPerSample
					totalSamples = recs[i].Detail.TotalSamples
					break
				case 4:
					vendorString = recs[i].Detail.VendorString
					comments     = strings.Join(recs[i].Detail.Comments,"\n")
					break

				case 6:
					ti := new(Image)
					ti.pictureType  = recs[i].Detail.PictureType
					ti.mimeType     = recs[i].Detail.MIMEType
					ti.description  = recs[i].Detail.Description
					ti.width        = recs[i].Detail.Width
					ti.height       = recs[i].Detail.Height
					ti.depth        = recs[i].Detail.Depth
					ti.dataLength   = recs[i].Detail.DataLength
					ti.imageData    = recs[i].Detail.ImageData
					images = append(images,*ti)
					break

				default:
					break
				
			}
		}

	
		/* check if record already exists */
		var sqlstr = `SELECT idx FROM tracks WHERE md5Signature=$1`
		st,dberr := db.Prepare(sqlstr)
		if dberr==nil {
			row,err := st.Query(md5Signature)
			var idx int
			row.Next()
			err = row.Scan(&idx)
			if err != nil {
				idx = 0
			}
			if (idx<1) {
			    sqlstr = `INSERT INTO tracks (
			      idx,
			      path,
			      md5Signature,
			      minBlockSize,
			      maxBlockSize,
			      minFrameSize,
			      maxFrameSize,
			      sampleRate,
			      channels,
			      bitsPerSample,
			      totalSamples,
			      vendorString,
			      comments,
			      sequence
			    ) VALUES (
			      DEFAULT,
			      $1,$2,$3,$4,$5,$6,
			      $7,$8,$9,$10,$11,$12,
			      CURRENT_TIMESTAMP
			    )`

			  st,dberr = db.Prepare(sqlstr)
			  if dberr != nil {
			  	log.Fatal(dberr)
			  }
			  st.Exec(path,
			    md5Signature,
			    minBlockSize,
			    maxBlockSize,
			    minFrameSize,
			    maxFrameSize,
			    sampleRate,
			    channels,
			    bitsPerSample,
			    totalSamples,
			    vendorString,
			    comments)

			  
			  /* insert images */
			  for i:= range images {
			    sqlstr = `INSERT INTO images (idx,md5Signature,
				      picturetype,mimetype,description,width,height,depth,
				      datalength,imagedata) VALUES (DEFAULT,$1,$2,$3,$4,
				      $5,$6,$7,$8,$9)`;
			    st,dberr = db.Prepare(sqlstr)
			    if dberr != nil {
				log.Fatal(dberr)
			    }
			    st.Exec(md5Signature,images[i].pictureType,images[i].mimeType,
				    images[i].description,images[i].width,images[i].height,
				    images[i].depth,images[i].dataLength,images[i].imageData)
			  }
			  fmt.Println(path + " complete.")
		      } else {
			  fmt.Println(path + " already in database.")
		      }
		 } else {
		    fmt.Println(fmt.Sprint(err))
		 }
	}
	return nil
} 

var db *sql.DB;

func main() {

  /* usage: goflac [path] 
   * path is optional, if specified it will scan the director and all subdirectories 
   * looking for *.flac files, parse the tags found using metaflac (modified for JSON
   * output) and insert them into the database.
   * before running first time, run the command
   * # createdb goflac
   * (presuming you have postgresql database server running on localhost and your 
   * user account has permission to create databases and tables)
   */
  var err error;
  db, err = sql.Open("postgres", "dbname=goflac sslmode=disable") 
  defer db.Close()

  if err != nil {
	log.Fatal(err)
  }
  _, err = db.Query("SELECT idx FROM tracks")
  if err != nil {
	/* tracks does not exist, try to create */
	_,err = db.Query(`CREATE TABLE tracks (
				idx serial,
				path character varying,
				md5Signature character varying, 
				minBlockSize integer, 
				maxBlockSize integer, 
				minFrameSize integer, 
				maxFrameSize integer, 
				sampleRate integer, 
				channels integer, 
				bitsPerSample integer, 
				totalSamples integer, 
				vendorString character varying, 
				comments text,
				sequence timestamp
			)`);
	if err != nil {
		log.Fatal(err);
	}
	_,err = db.Query(`CREATE TABLE images (
				idx serial,
				blockid integer,
				md5Signature character varying,
				picturetype integer,
				mimetype character varying,
				description character varying,
				width integer,
				height integer,
				depth integer,
				datalength integer,
				imagedata text
			)`);
	if err != nil {
		log.Fatal(err);
	}
  }

  flag.Parse()
  root := flag.Arg(0)
  if (root != "") {
	err := filepath.Walk(root, visit)
  	if (err != nil) {
  		fmt.Printf("filepath.Walk() returned %v\n", err)
  	}
  }
  
  fmt.Println()
  fmt.Println("#\tChan\tRate\tBPS\tPath")
  fmt.Println("---\t----\t------\t---\t----")
  var m map[int]string
  m = make(map[int]string)
  var d map[int]string
  d = make(map[int]string)
  
  var sqlstr = `SELECT idx,path,channels,samplerate,bitspersample,comments FROM tracks ORDER BY sequence DESC`
  st,err := db.Prepare(sqlstr)
  row,err := st.Query()
  if err==nil {
    for row.Next() {
      var idx int
      var path string
      var channels int
      var samplerate int
      var bitspersample int
      var comments string
      row.Scan(&idx,&path,&channels,&samplerate,&bitspersample,&comments)
      fmt.Printf("%d\t%d\t%d\t%d\t%s\n",idx,channels,samplerate,bitspersample,path);
      m[idx]=path;
      d[idx]=comments;
    }
  }
  
  fmt.Println()

  var selected_idx int
  for ;; {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter the Track Number: (Control-C to exit) \n")
    u, _ := reader.ReadString('\n')
    selected_idx,err = strconv.Atoi(strings.TrimSpace(u))
    if err != nil {
      fmt.Println("Invalid ID")
    } else {
      if selected_idx<1 {
	  fmt.Println("Invalid ID")
      } else {
	/* play the tune */
	if m[selected_idx]=="" {
	  fmt.Println("Invalid ID")
	} else {
	  fmt.Println()
	  fmt.Println("Now Playing: " + m[selected_idx])
	  fmt.Println()

	  dx := strings.Split(d[selected_idx],"\n")
	  for i:= range dx {
		rx := strings.Split(dx[i],"=") /* hmmm */
		fmt.Printf("%-40s\t%s\n",rx[0],rx[1])
          }	
	  fmt.Println()
	  fmt.Println(" (Press Return to Quit)")
	  fmt.Println()
	  mpargs := []string{m[selected_idx]};
	  cmd := exec.Command("/usr/local/bin/mplayer",mpargs...);
	  cmd.Start()
	  _,_ = reader.ReadString('\n')
	  err = cmd.Process.Kill() /* oh, no */
	}
      }
    }
  }
 
}
