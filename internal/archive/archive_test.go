package archive

import (
 "archive/tar"
 "archive/zip"
 "bytes"
 "os"
 "path/filepath"
 "testing"
)
func makeZip(t *testing.T,path string, files map[string]string){t.Helper();f,e:=os.Create(path);if e!=nil{t.Fatal(e)};z:=zip.NewWriter(f);for n,v:=range files{w,e:=z.Create(n);if e!=nil{t.Fatal(e)};if _,e=w.Write([]byte(v));e!=nil{t.Fatal(e)}};if e=z.Close();e!=nil{t.Fatal(e)};if e=f.Close();e!=nil{t.Fatal(e)}}
func TestInspectAndExtractZip(t *testing.T){d:=t.TempDir();p:=filepath.Join(d,"x.zip");makeZip(t,p,map[string]string{"docs/a.txt":"hello","b.txt":"world"});s,e:=Inspect(p);if e!=nil{t.Fatal(e)};if s.Format!="zip"||s.Files!=2||s.UncompressedBytes!=10{t.Fatalf("unexpected summary: %+v",s)};out:=filepath.Join(d,"out");if e=Extract(p,out,false);e!=nil{t.Fatal(e)};b,e:=os.ReadFile(filepath.Join(out,"docs","a.txt"));if e!=nil||string(b)!="hello"{t.Fatalf("bad extraction: %q %v",b,e)}}
func TestRejectZipSlip(t *testing.T){d:=t.TempDir();p:=filepath.Join(d,"bad.zip");makeZip(t,p,map[string]string{"../escape.txt":"no"});if e:=Extract(p,filepath.Join(d,"out"),false);e==nil{t.Fatal("expected unsafe path rejection")};if _,e:=os.Stat(filepath.Join(d,"escape.txt"));!os.IsNotExist(e){t.Fatal("archive escaped destination")}}
func TestRefuseOverwrite(t *testing.T){d:=t.TempDir();p:=filepath.Join(d,"x.zip");makeZip(t,p,map[string]string{"a.txt":"new"});out:=filepath.Join(d,"out");if e:=os.Mkdir(out,0755);e!=nil{t.Fatal(e)};if e:=os.WriteFile(filepath.Join(out,"a.txt"),[]byte("old"),0644);e!=nil{t.Fatal(e)};if e:=Extract(p,out,false);e==nil{t.Fatal("expected overwrite refusal")}}
func TestTarTraversalRejected(t *testing.T){d:=t.TempDir();p:=filepath.Join(d,"x.tar");f,e:=os.Create(p);if e!=nil{t.Fatal(e)};tw:=tar.NewWriter(f);data:=[]byte("x");if e=tw.WriteHeader(&tar.Header{Name:"../../bad",Mode:0644,Size:int64(len(data))});e!=nil{t.Fatal(e)};_,_=tw.Write(data);_=tw.Close();_=f.Close();if e=Extract(p,filepath.Join(d,"out"),false);e==nil{t.Fatal("expected traversal rejection")}}
func TestSafePath(t *testing.T){root:=t.TempDir();p,e:=safePath(root,"a/b.txt");if e!=nil||p!=filepath.Join(root,"a","b.txt"){t.Fatalf("unexpected %q %v",p,e)};for _,n:=range []string{"../x","a/../../x"}{if _,e=safePath(root,n);e==nil{t.Fatalf("accepted %q",n)}}}
func TestTarInspect(t *testing.T){d:=t.TempDir();p:=filepath.Join(d,"x.tar");f,_:=os.Create(p);tw:=tar.NewWriter(f);data:=[]byte("abc");_ = tw.WriteHeader(&tar.Header{Name:"a.txt",Mode:0644,Size:3});_,_=tw.Write(data);_=tw.Close();_=f.Close();s,e:=Inspect(p);if e!=nil{t.Fatal(e)};if s.Format!="tar"||s.Files!=1||s.UncompressedBytes!=3{t.Fatalf("bad summary %+v",s)}}
func TestUnsupported(t *testing.T){p:=filepath.Join(t.TempDir(),"x.bin");_ = os.WriteFile(p,bytes.Repeat([]byte{1},600),0644);if _,e:=Inspect(p);e==nil{t.Fatal("expected unsupported error")}}
