# Archive Explorer

A small, dependency-free Go CLI for inspecting and safely extracting **ZIP, TAR, and TAR.GZ/TGZ** archives. It is designed for local workflows where predictable behavior and extraction safety matter more than a large feature surface.

## Why it exists
Archive files often need a quick inspection before extraction. Archive Explorer lists contents without unpacking them, emits machine-readable JSON, and extracts with explicit overwrite behavior and defenses against path traversal and symbolic-link entries.

## Features
- Detects ZIP, TAR, and gzip-compressed TAR from file content rather than filename alone.
- Lists entries, type, uncompressed size, totals, and archive format.
- `--json` output for scripts and automation.
- Safe extraction with Zip Slip/path-traversal checks.
- Rejects symlink and special TAR entries instead of materializing them.
- Refuses to overwrite existing files unless `--overwrite` is explicitly supplied.
- Uses temporary files followed by rename when writing extracted regular files.
- No network access, telemetry, configuration file, API key, or runtime dependency.

## Requirements
- Go 1.22+ to build from source.

## Installation
```bash
git clone https://github.com/rad03i2/archive-explorer.git
cd archive-explorer
go build -o archive-explorer ./cmd/archive-explorer
```

Or run without installing:
```bash
go run ./cmd/archive-explorer list ./backup.zip
```

## Usage
```bash
# Human-readable listing
./archive-explorer list backup.zip

# Machine-readable listing
./archive-explorer list --json backup.tar.gz

# Safe extraction; existing files are protected
./archive-explorer extract backup.zip ./output

# Explicitly allow replacing existing regular files
./archive-explorer extract --overwrite backup.zip ./output

./archive-explorer version
```

### Preview guidance
This is a terminal application, so screenshots are optional. For a project preview, capture the `list` table and a JSON listing rather than implying a graphical interface.

## Configuration
There is no persistent configuration. All behavior is explicit through CLI arguments. `--overwrite` is intentionally opt-in.

## Project structure
```text
cmd/archive-explorer/main.go       CLI and output formatting
internal/archive/archive.go        detection, inspection, safe extraction
internal/archive/archive_test.go   functional and security tests
.github/workflows/ci.yml           cross-platform build/test pipeline
```

## Testing
```bash
gofmt -w .
go vet ./...
go test ./... -race
go build ./cmd/archive-explorer
```
Tests cover ZIP/TAR inspection, extraction, overwrite protection, unsupported input, and traversal attempts.

## Security & privacy
All processing is local. Extraction validates every destination path against the requested output root. Link/special entries are rejected because safely reproducing their semantics is platform-dependent. See [SECURITY.md](SECURITY.md) for reporting guidance.

## Limitations
- Supports ZIP, TAR, and TAR.GZ/TGZ only; 7z/RAR are intentionally not implemented.
- Does not decrypt password-protected archives.
- Does not create archives or edit archive contents.
- Metadata such as ownership is not restored; regular file permission bits are applied where the OS supports them.
- Extraction is fail-fast, not transactional: files successfully written before a later invalid entry are not automatically rolled back.

## Optional roadmap
Potential future work includes selective extraction and additional formats through carefully reviewed libraries. These are not current features.

## Contributing
See [CONTRIBUTING.md](CONTRIBUTING.md). Please include tests for behavioral or security-sensitive changes.

## License
MIT — see [LICENSE](LICENSE).

## Author
**Radwan Abdulhadi Ahmed**  
**رضوان عبدالهادي أحمد**  
GitHub: **@rad03i2**

---

# العربية — Archive Explorer

أداة سطر أوامر صغيرة مكتوبة بلغة **Go** وبدون مكتبات تشغيل خارجية، لفحص واستخراج أرشيفات **ZIP وTAR وTAR.GZ/TGZ** محليًا بصورة آمنة ومتوقعة.

## لماذا هذا المشروع؟
نحتاج كثيرًا إلى معرفة محتويات الأرشيف قبل فكّه. يعرض Archive Explorer الملفات دون استخراجها، ويوفر JSON للأتمتة، ويجعل استبدال الملفات خيارًا صريحًا بدل أن يحدث تلقائيًا، مع حماية من محاولات الخروج خارج مجلد الاستخراج.

## المميزات
- اكتشاف نوع الأرشيف من محتوى الملف وليس الامتداد فقط.
- عرض الملفات والمجلدات والأحجام والإجماليات ونوع الأرشيف.
- إخراج JSON للاستخدام البرمجي.
- حماية من Zip Slip ومسارات `../` الخطرة.
- رفض الروابط الرمزية والعناصر الخاصة داخل TAR/ZIP أثناء الاستخراج.
- عدم استبدال ملف موجود إلا عند تمرير `--overwrite` صراحةً.
- كتابة الملفات عبر ملف مؤقت ثم إعادة تسميته.
- لا اتصال بالإنترنت ولا Telemetry ولا مفاتيح API ولا إعدادات سرية.

## المتطلبات والتثبيت
يتطلب Go 1.22 أو أحدث:
```bash
git clone https://github.com/rad03i2/archive-explorer.git
cd archive-explorer
go build -o archive-explorer ./cmd/archive-explorer
```

## الاستخدام
```bash
./archive-explorer list backup.zip
./archive-explorer list --json backup.tar.gz
./archive-explorer extract backup.zip ./output
./archive-explorer extract --overwrite backup.zip ./output
./archive-explorer version
```

### المعاينة والصور
المشروع أداة طرفية وليس واجهة رسومية. عند إضافة صورة تعريفية للمستودع يُفضّل تصوير جدول `list` أو مخرجات JSON الحقيقية.

## الإعداد
لا يوجد ملف إعداد دائم. الخيارات كلها صريحة في سطر الأوامر، وميزة الاستبدال معطلة افتراضيًا.

## هيكل المشروع
- `cmd/archive-explorer/main.go`: واجهة CLI.
- `internal/archive/archive.go`: الكشف والفحص والاستخراج الآمن.
- `internal/archive/archive_test.go`: اختبارات وظيفية وأمنية.
- `.github/workflows/ci.yml`: اختبارات وبناء متعدد الأنظمة.

## الاختبارات
```bash
gofmt -w .
go vet ./...
go test ./... -race
go build ./cmd/archive-explorer
```
تشمل الاختبارات ZIP وTAR، والاستخراج، ومنع الاستبدال، والمدخلات غير المدعومة، ومحاولات تجاوز مسار الوجهة.

## الأمان والخصوصية
كل المعالجة محلية. يتم التحقق من كل مسار قبل الكتابة داخل مجلد الوجهة، وتُرفض الروابط والعناصر الخاصة بدل التعامل معها بطريقة قد تختلف أمنيًا بين الأنظمة. راجع [SECURITY.md](SECURITY.md).

## القيود
- لا يدعم RAR أو 7z حاليًا.
- لا يفك تشفير الأرشيفات المحمية بكلمة مرور.
- لا ينشئ الأرشيفات ولا يعدّل محتواها.
- لا يستعيد ملكية الملفات، وتطبيق صلاحيات الملفات يعتمد على نظام التشغيل.
- الاستخراج يتوقف عند أول خطأ وليس معاملة ذرية كاملة؛ الملفات التي استُخرجت قبل الخطأ لا تُحذف تلقائيًا.

## تطويرات مستقبلية اختيارية
يمكن مستقبلًا إضافة استخراج ملفات محددة أو صيغ إضافية بعد مراجعة المكتبات المناسبة. هذه ليست ميزات حالية.

## المساهمة والترخيص
راجع [CONTRIBUTING.md](CONTRIBUTING.md). المشروع مرخص وفق MIT في [LICENSE](LICENSE).

## المؤلف
**Radwan Abdulhadi Ahmed**  
**رضوان عبدالهادي أحمد**  
GitHub: **@rad03i2**
