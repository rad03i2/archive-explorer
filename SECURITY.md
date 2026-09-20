# Security Policy / سياسة الأمان

## Supported version
The current `main` branch is supported.

## Reporting
Please report security concerns privately through GitHub's available private security-reporting mechanism when enabled. Do not publish exploit archives or sensitive details in a public issue. Include the affected command, archive format, operating system, minimal reproduction steps, and expected/actual behavior. Never include credentials or personal data.

Archive Explorer treats archive paths and link entries as untrusted. The extractor rejects traversal, symbolic links, and special TAR entries and protects existing files unless overwrite is explicit. These controls reduce risk but do not make arbitrary archives trustworthy.

## العربية
يتم دعم النسخة الحالية من فرع `main`. عند اكتشاف مشكلة أمنية، استخدم وسيلة الإبلاغ الأمني الخاص في GitHub إن كانت مفعلة، ولا تنشر ملف استغلال أو تفاصيل حساسة في Issue عامة. اذكر الأمر والصيغة ونظام التشغيل وخطوات إعادة المشكلة دون بيانات شخصية أو أسرار.

يتعامل المشروع مع مسارات الأرشيف والروابط كمدخلات غير موثوقة، ويرفض تجاوز المسار والروابط الرمزية والعناصر الخاصة ويحمي الملفات الموجودة افتراضيًا. هذه الحمايات تقلل المخاطر لكنها لا تجعل كل أرشيف مجهول موثوقًا.
