package main

import (
 "encoding/json"
 "flag"
 "fmt"
 "os"
 "text/tabwriter"

 arc "github.com/rad03i2/archive-explorer/internal/archive"
)

const version = "1.0.0"
func usage(){fmt.Fprintln(os.Stderr,"Archive Explorer — inspect and safely extract ZIP/TAR/TAR.GZ archives\n\nUsage:\n  archive-explorer list [--json] ARCHIVE\n  archive-explorer extract [--overwrite] ARCHIVE DESTINATION\n  archive-explorer version")}
func main(){if len(os.Args)<2{usage();os.Exit(2)};var err error;switch os.Args[1]{case "list":err=list(os.Args[2:]);case "extract":err=extract(os.Args[2:]);case "version","--version","-v":fmt.Printf("archive-explorer %s — Radwan Abdulhadi Ahmed / @rad03i2\n",version);return;case "help","--help","-h":usage();return;default:usage();os.Exit(2)};if err!=nil{fmt.Fprintln(os.Stderr,"error:",err);os.Exit(1)}}
func list(args []string)error{fs:=flag.NewFlagSet("list",flag.ContinueOnError);asJSON:=fs.Bool("json",false,"emit JSON");if err:=fs.Parse(args);err!=nil{return err};if fs.NArg()!=1{return fmt.Errorf("list requires exactly one archive path")};s,err:=arc.Inspect(fs.Arg(0));if err!=nil{return err};if *asJSON{b,e:=json.MarshalIndent(s,"","  ");if e!=nil{return e};fmt.Println(string(b));return nil};fmt.Printf("Archive: %s\nFormat: %s | Files: %d | Directories: %d | Uncompressed: %d bytes\n\n",s.Path,s.Format,s.Files,s.Directories,s.UncompressedBytes);w:=tabwriter.NewWriter(os.Stdout,0,4,2,' ',0);fmt.Fprintln(w,"TYPE\tSIZE\tNAME");for _,e:=range s.Entries{kind:="file";if e.Directory{kind="dir"};fmt.Fprintf(w,"%s\t%d\t%s\n",kind,e.Size,e.Name)};return w.Flush()}
func extract(args []string)error{fs:=flag.NewFlagSet("extract",flag.ContinueOnError);overwrite:=fs.Bool("overwrite",false,"replace existing regular files");if err:=fs.Parse(args);err!=nil{return err};if fs.NArg()!=2{return fmt.Errorf("extract requires ARCHIVE and DESTINATION")};if err:=arc.Extract(fs.Arg(0),fs.Arg(1),*overwrite);err!=nil{return err};fmt.Println("Extraction completed safely.");return nil}
