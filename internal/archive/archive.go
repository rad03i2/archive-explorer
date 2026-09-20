package archive

import (
 "archive/tar"
 "archive/zip"
 "compress/gzip"
 "errors"
 "fmt"
 "io"
 "os"
 "path/filepath"
 "sort"
 "strings"
 "time"
)

type Entry struct { Name string `json:"name"`; Size int64 `json:"size"`; Mode string `json:"mode"`; Modified time.Time `json:"modified"`; Directory bool `json:"directory"` }
type Summary struct { Path string `json:"path"`; Format string `json:"format"`; Entries []Entry `json:"entries"`; Files int `json:"files"`; Directories int `json:"directories"`; UncompressedBytes int64 `json:"uncompressed_bytes"` }

func Detect(path string) (string, error) {
 f, err := os.Open(path); if err != nil { return "", err }; defer f.Close()
 h := make([]byte, 512); n, _ := io.ReadFull(f, h); h = h[:n]
 if n >= 4 && string(h[:4]) == "PK\x03\x04" { return "zip", nil }
 if n >= 2 && h[0] == 0x1f && h[1] == 0x8b { return "tar.gz", nil }
 if n >= 262 && string(h[257:262]) == "ustar" { return "tar", nil }
 return "", errors.New("unsupported archive: supported formats are ZIP, TAR, TAR.GZ/TGZ")
}

func Inspect(path string) (Summary, error) {
 abs, err := filepath.Abs(path); if err != nil { return Summary{}, err }
 st, err := os.Stat(abs); if err != nil { return Summary{}, err }; if !st.Mode().IsRegular() { return Summary{}, errors.New("archive path must be a regular file") }
 format, err := Detect(abs); if err != nil { return Summary{}, err }
 s := Summary{Path: abs, Format: format}
 switch format { case "zip": err = inspectZip(abs, &s); case "tar": err = inspectTar(abs, false, &s); case "tar.gz": err = inspectTar(abs, true, &s) }
 if err != nil { return Summary{}, err }
 sort.Slice(s.Entries, func(i,j int) bool { return s.Entries[i].Name < s.Entries[j].Name })
 for _, e := range s.Entries { if e.Directory { s.Directories++ } else { s.Files++; s.UncompressedBytes += e.Size } }
 return s, nil
}
func inspectZip(path string, s *Summary) error { r, err := zip.OpenReader(path); if err != nil { return fmt.Errorf("open zip: %w", err) }; defer r.Close(); for _, f := range r.File { s.Entries=append(s.Entries, Entry{f.Name,int64(f.UncompressedSize64),f.Mode().String(),f.Modified,f.FileInfo().IsDir()}) }; return nil }
func inspectTar(path string, gz bool, s *Summary) error { f,err:=os.Open(path); if err!=nil{return err}; defer f.Close(); var r io.Reader=f; if gz { g,e:=gzip.NewReader(f); if e!=nil{return fmt.Errorf("open gzip: %w",e)}; defer g.Close(); r=g }; tr:=tar.NewReader(r); for { h,e:=tr.Next(); if e==io.EOF{break}; if e!=nil{return fmt.Errorf("read tar: %w",e)}; s.Entries=append(s.Entries,Entry{h.Name,h.Size,os.FileMode(h.Mode).String(),h.ModTime,h.FileInfo().IsDir()}) }; return nil }

func Extract(path, dest string, overwrite bool) error {
 format, err := Detect(path); if err != nil{return err}; root,err:=filepath.Abs(dest); if err!=nil{return err}; if err=os.MkdirAll(root,0755);err!=nil{return err}
 switch format { case "zip": return extractZip(path,root,overwrite); case "tar": return extractTar(path,root,false,overwrite); case "tar.gz": return extractTar(path,root,true,overwrite) }; return nil
}
func safePath(root,name string)(string,error){ clean:=filepath.Clean(filepath.FromSlash(name)); if clean=="." || filepath.IsAbs(clean) || clean==".." || strings.HasPrefix(clean,".."+string(os.PathSeparator)){return "",fmt.Errorf("unsafe archive path %q",name)}; out:=filepath.Join(root,clean); rel,e:=filepath.Rel(root,out); if e!=nil||rel==".."||strings.HasPrefix(rel,".."+string(os.PathSeparator)){return "",fmt.Errorf("path escapes destination: %q",name)}; return out,nil }
func writeFile(out string, mode os.FileMode, overwrite bool, src io.Reader) error { if st,e:=os.Lstat(out); e==nil { if st.Mode()&os.ModeSymlink!=0{return fmt.Errorf("refusing symlink destination %q",out)}; if !overwrite{return fmt.Errorf("destination exists: %q",out)} } else if !os.IsNotExist(e){return e}; if err:=os.MkdirAll(filepath.Dir(out),0755);err!=nil{return err}; tmp,err:=os.CreateTemp(filepath.Dir(out),".archive-explorer-*");if err!=nil{return err}; tmpName:=tmp.Name(); defer os.Remove(tmpName); if _,err=io.Copy(tmp,src);err!=nil{tmp.Close();return err}; if err=tmp.Close();err!=nil{return err}; if err=os.Chmod(tmpName,mode.Perm());err!=nil{return err}; if overwrite{_ = os.Remove(out)}; return os.Rename(tmpName,out) }
func extractZip(path,root string,overwrite bool)error{ r,e:=zip.OpenReader(path);if e!=nil{return e};defer r.Close(); for _,f:=range r.File{ out,e:=safePath(root,f.Name);if e!=nil{return e}; if f.Mode()&os.ModeSymlink!=0{return fmt.Errorf("refusing symbolic link %q",f.Name)}; if f.FileInfo().IsDir(){if e=os.MkdirAll(out,f.Mode().Perm());e!=nil{return e};continue}; rc,e:=f.Open();if e!=nil{return e}; e=writeFile(out,f.Mode(),overwrite,rc);rc.Close();if e!=nil{return e} };return nil }
func extractTar(path,root string,gz,overwrite bool)error{ f,e:=os.Open(path);if e!=nil{return e};defer f.Close();var r io.Reader=f;if gz{g,x:=gzip.NewReader(f);if x!=nil{return x};defer g.Close();r=g};tr:=tar.NewReader(r);for{h,x:=tr.Next();if x==io.EOF{break};if x!=nil{return x};out,x:=safePath(root,h.Name);if x!=nil{return x};switch h.Typeflag{case tar.TypeDir:x=os.MkdirAll(out,os.FileMode(h.Mode).Perm());case tar.TypeReg,tar.TypeRegA:x=writeFile(out,os.FileMode(h.Mode),overwrite,tr);default:return fmt.Errorf("refusing unsupported link/special entry %q",h.Name)};if x!=nil{return x}};return nil}
