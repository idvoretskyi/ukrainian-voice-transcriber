# Ukrainian Voice Transcriber - Посібник користувача

## Зміст
1. [Вступ](#вступ)
2. [Встановлення](#встановлення)
3. [Налаштування автентифікації](#налаштування-автентифікації)
4. [Основне використання](#основне-використання)
5. [Розширені можливості](#розширені-можливості)
6. [Організація файлів](#організація-файлів)
7. [Виправлення проблем](#виправлення-проблем)
8. [Управління витратами](#управління-витратами)
9. [Поширені запитання](#поширені-запитання)

## Вступ

Ukrainian Voice Transcriber - це інструмент командного рядка, який конвертує українські аудіо та відео файли в текст за допомогою Google Cloud Speech-to-Text API. Він розроблений для контент-мейкерів, журналістів, дослідників та всіх, хто потребує транскрибування українського контенту.

### Основні можливості
- **Оптимізовано для української мови**: Спеціально налаштовано для розпізнавання української мови (`uk-UA`)
- **Множинні формати**: Підтримує всі відео формати, які може обробляти FFmpeg (MP4, AVI, MOV, MKV тощо)
- **Розумна організація файлів**: Автоматично створює впорядковані директорії для транскриптів
- **Економічно ефективний**: Використовує стандартні (не преміум) моделі Speech-to-Text з автоматичним очищенням
- **Просте налаштування**: Працює з gcloud CLI - не потребує складної конфігурації OAuth

### Системні вимоги
- **Операційна система**: macOS, Linux або Windows
- **FFmpeg**: Необхідний для вилучення аудіо з відео файлів
- **Google Cloud проект**: З увімкненими API Speech-to-Text та Cloud Storage
- **Інтернет-з'єднання**: Необхідне для викликів API

## Встановлення

### Крок 1: Встановлення FFmpeg

**macOS:**
```bash
brew install ffmpeg
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install ffmpeg
```

**Windows:**
1. Завантажте FFmpeg з https://ffmpeg.org/download.html
2. Розпакуйте у теку (наприклад, `C:\ffmpeg`)
3. Додайте `C:\ffmpeg\bin` до змінної середовища PATH

**Перевірка встановлення:**
```bash
ffmpeg -version
```

### Крок 2: Встановлення Ukrainian Voice Transcriber

**Варіант A: Завантаження готового файлу**
1. Перейдіть на [сторінку релізів](https://github.com/idvoretskyi/ukrainian-voice-transcriber/releases)
2. Завантажте файл для вашої операційної системи
3. Зробіть його виконуваним (Linux/macOS):
   ```bash
   chmod +x ukrainian-voice-transcriber
   ```

**Варіант B: Збірка з коду**
```bash
# Встановіть Go (1.24 або новіший)
# Завантажте з https://golang.org/dl/

# Клонуйте та зберіть
git clone https://github.com/idvoretskyi/ukrainian-voice-transcriber.git
cd ukrainian-voice-transcriber
go build -o ukrainian-voice-transcriber cmd/transcriber/main.go
```

### Крок 3: Перевірка встановлення
```bash
./ukrainian-voice-transcriber --help
```

## Налаштування автентифікації

Додаток підтримує два методи автентифікації. Ми рекомендуємо використовувати gcloud CLI для простоти.

### Метод 1: gcloud CLI (Рекомендовано)

**Крок 1: Встановлення gcloud CLI**
```bash
# macOS
brew install google-cloud-sdk

# Linux
curl https://sdk.cloud.google.com | bash
exec -l $SHELL

# Windows - Завантажте інсталятор з:
# https://cloud.google.com/sdk/docs/install
```

**Крок 2: Автентифікація та налаштування**
```bash
# Увійдіть у свій Google акаунт
gcloud auth login

# Налаштуйте стандартні облікові дані для додатків
gcloud auth application-default login

# Встановіть ID вашого проекту (замініть на ваш реальний ID проекту)
gcloud config set project ВАШ_ID_ПРОЕКТУ

# Увімкніть необхідні API
gcloud services enable speech.googleapis.com
gcloud services enable storage.googleapis.com
```

**Крок 3: Перевірка налаштувань**
```bash
# Перевірте статус автентифікації
./ukrainian-voice-transcriber auth --status

# Запустіть перевірку налаштувань
./ukrainian-voice-transcriber setup
```

### Метод 2: Сервісний акаунт (Для досвідчених)

**Крок 1: Створення сервісного акаунта**
1. Перейдіть до [Google Cloud Console](https://console.cloud.google.com/)
2. Виберіть ваш проект або створіть новий
3. Перейдіть до "IAM & Admin" → "Service Accounts"
4. Натисніть "Create Service Account"
5. Введіть назву (наприклад, "ukrainian-transcriber")
6. Натисніть "Create and Continue"

**Крок 2: Призначення ролей**
Додайте ці ролі:
- `Cloud Speech Client`
- `Storage Admin`

**Крок 3: Створення ключа**
1. Натисніть на створений сервісний акаунт
2. Перейдіть до вкладки "Keys"
3. Натисніть "Add Key" → "Create New Key"
4. Виберіть формат "JSON"
5. Завантажте файл ключа

**Крок 4: Налаштування**
```bash
# Помістіть файл ключа у директорію додатка
cp ~/Downloads/service-account-key.json ./service-account.json

# Або встановіть змінну середовища
export GOOGLE_APPLICATION_CREDENTIALS="/шлях/до/service-account.json"
```

## Основне використання

### Перше налаштування
```bash
# Перевірте, чи все налаштовано правильно
./ukrainian-voice-transcriber setup
```

### Базове транскрибування
```bash
# Транскрибуйте відео файл
./ukrainian-voice-transcriber transcribe video.mp4
```

Це зробить:
1. Створить директорію з назвою `video/`
2. Збереже транскрипт як `video/video.txt`
3. Покаже прогрес та результати

### Відео з пробілами в назві
```bash
# Автоматично обробляє пробіли
./ukrainian-voice-transcriber transcribe "Моє Інтерв'ю 2024.mp4"
```

Це зробить:
1. Створить директорію з назвою `Моє_Інтерв'ю_2024/`
2. Збереже транскрипт як `Моє_Інтерв'ю_2024/Моє_Інтерв'ю_2024.txt`

### Збереження у конкретний файл
```bash
# Перевизначити стандартну організацію файлів
./ukrainian-voice-transcriber transcribe video.mp4 -o власний_транскрипт.txt
```

### Детальний вивід
```bash
# Переглянути детальну інформацію про обробку
./ukrainian-voice-transcriber transcribe video.mp4 --verbose
```

### Тихий режим
```bash
# Показати тільки фінальний транскрипт
./ukrainian-voice-transcriber transcribe video.mp4 --quiet
```

## Розширені можливості

### Пакетна обробка
```bash
# Обробити всі MP4 файли в поточній директорії
for video in *.mp4; do
    echo "Обробляється: $video"
    ./ukrainian-voice-transcriber transcribe "$video"
done
```

### Власне сховище
```bash
# Використовувати власний Cloud Storage bucket
./ukrainian-voice-transcriber transcribe video.mp4 --bucket мій-власний-bucket
```

### Інтеграція зі скриптами
```bash
#!/bin/bash
# Приклад: Обробка та підрахунок слів

VIDEO_FILE="$1"
if [ -z "$VIDEO_FILE" ]; then
    echo "Використання: $0 <відео-файл>"
    exit 1
fi

echo "Транскрибування $VIDEO_FILE..."
./ukrainian-voice-transcriber transcribe "$VIDEO_FILE" --quiet > транскрипт.txt

WORD_COUNT=$(wc -w < транскрипт.txt)
echo "Транскрибування завершено. Кількість слів: $WORD_COUNT"
```

## Організація файлів

### Стандартна організація
Коли ви запускаєте:
```bash
./ukrainian-voice-transcriber transcribe "Моє Відео.mp4"
```

Додаток створює:
```
Моє_Відео/
└── Моє_Відео.txt
```

### Кілька файлів
```bash
./ukrainian-voice-transcriber transcribe інтерв'ю1.mp4
./ukrainian-voice-transcriber transcribe інтерв'ю2.mp4
./ukrainian-voice-transcriber transcribe "Фінальна Дискусія.mp4"
```

Результат:
```
інтерв'ю1/
└── інтерв'ю1.txt
інтерв'ю2/
└── інтерв'ю2.txt
Фінальна_Дискусія/
└── Фінальна_Дискусія.txt
```

### Власний вивід
```bash
# Збереження у конкретне місце
./ukrainian-voice-transcriber transcribe video.mp4 -o транскрипти/мій-транскрипт.txt
```

## Виправлення проблем

### Поширені проблеми

**1. "FFmpeg not found"**
```bash
# Перевірте, чи встановлено FFmpeg
ffmpeg -version

# Встановіть, якщо відсутній
brew install ffmpeg  # macOS
sudo apt install ffmpeg  # Ubuntu
```

**2. "Authentication required"**
```bash
# Перевірте статус автентифікації
./ukrainian-voice-transcriber auth --status

# Якщо використовуєте gcloud
gcloud auth login
gcloud auth application-default login

# Якщо використовуєте сервісний акаунт
ls -la service-account.json
```

**3. "Permission denied"**
```bash
# Зробіть файл виконуваним
chmod +x ukrainian-voice-transcriber

# Або запустіть з повним шляхом
./ukrainian-voice-transcriber transcribe video.mp4
```

**4. "Project not found" або "Unknown project id"**
```bash
# Встановіть ID проекту
gcloud config set project ВАШ_ID_ПРОЕКТУ

# Перевірте проект
gcloud config get-value project
```

**5. "API not enabled"**
```bash
# Увімкніть необхідні API
gcloud services enable speech.googleapis.com
gcloud services enable storage.googleapis.com
```

### Режим налагодження
```bash
# Отримати детальну інформацію про помилки
./ukrainian-voice-transcriber transcribe video.mp4 --verbose
```

### Файли логів
Додаток не створює файли логів за замовчуванням. Для налагодження:
```bash
# Перенаправити вивід у файл
./ukrainian-voice-transcriber transcribe video.mp4 --verbose > налагодження.log 2>&1
```

## Управління витратами

### Розуміння витрат
- **Speech-to-Text**: ~$0.006 за 15-секундний фрагмент
- **Cloud Storage**: ~$0.020 за ГБ/місяць (тимчасові файли)
- **Типове 1-годинне відео**: ~$1.44

### Поради з оптимізації витрат

**1. Використання стандартної моделі**
Додаток за замовчуванням використовує стандартну (не покращену) модель Speech-to-Text для економії.

**2. Автоматичне очищення**
Тимчасові файли автоматично видаляються після обробки та мають політику життєвого циклу в 1 день.

**3. Моніторинг використання**
```bash
# Перевірте білінг Google Cloud
gcloud billing accounts list
gcloud billing budgets list
```

**4. Налаштування бюджетних сповіщень**
1. Перейдіть до [Google Cloud Console](https://console.cloud.google.com/)
2. Перейдіть до "Billing" → "Budgets & alerts"
3. Створіть бюджет з email сповіщеннями

### Оцінка витрат
```bash
# Для планування
echo "Тривалість: 60 хвилин"
echo "Оцінка витрат: $1.44 (60 хв × $0.006/15-сек × 4 фрагменти/хв)"
```

## Поширені запитання

### Загальні питання

**П: Які відео формати підтримуються?**
В: Всі формати, які підтримує FFmpeg (MP4, AVI, MOV, MKV, WebM тощо)

**П: Які аудіо формати підтримуються?**
В: MP3, WAV, FLAC, M4A та інші, які підтримує FFmpeg

**П: Чи можу я транскрибувати тільки аудіо файли?**
В: Так, додаток працює як з відео, так і з аудіо файлами

**П: Наскільки точне транскрибування?**
В: Точність залежить від якості аудіо, чіткості мовлення та фонового шуму. Розпізнавання української мови оптимізовано для локалі `uk-UA`.

**П: Чи можу я транскрибувати інші мови?**
В: Додаток спеціально оптимізований для української мови. Для інших мов потрібно буде змінити налаштування мови в коді.

### Технічні питання

**П: Де зберігаються тимчасові файли?**
В: У Google Cloud Storage в bucket з назвою `{project-id}-ukr-voice-transcriber`

**П: Як довго зберігаються тимчасові файли?**
В: Тимчасові файли видаляються одразу після обробки та мають політику життєвого циклу в 1 день як резерв

**П: Чи можу я використовувати власний storage bucket?**
В: Так, використовуйте прапорець `--bucket` для вказання власної назви bucket

**П: Що відбувається, якщо транскрибування не вдається?**
В: Тимчасові файли автоматично очищуються, а детальні повідомлення про помилки надаються

### Питання безпеки

**П: Чи безпечні мої дані?**
В: Так, файли обробляються через безпечну інфраструктуру Google Cloud, а тимчасові файли автоматично видаляються

**П: Чи зберігаються облікові дані локально?**
В: Використовуються тільки стандартні облікові дані додатків (через gcloud) або ключі сервісних акаунтів. Власні облікові дані не зберігаються

**П: Чи можу я використовувати це у продакшені?**
В: Так, але переконайтеся, що у вас є належна автентифікація, моніторинг та обробка помилок

### Питання продуктивності

**П: Скільки часу займає транскрибування?**
В: Час обробки варіюється, але зазвичай це 1-2x тривалості аудіо (наприклад, 10 хвилин для 5-хвилинного відео)

**П: Чи можу я обробляти кілька файлів одночасно?**
В: Додаток обробляє один файл за раз, але ви можете запустити кілька екземплярів паралельно

**П: Які обмеження розміру файлів?**
В: Google Cloud Speech-to-Text має обмеження в 10МБ для синхронних запитів, але додаток обробляє довші файли шляхом завантаження до Cloud Storage

---

**Потрібна додаткова допомога?**
- Перевірте [розділ виправлення проблем](#виправлення-проблем)
- Відвідайте [GitHub репозиторій](https://github.com/idvoretskyi/ukrainian-voice-transcriber)
- Переглянуйте [README.md](../README.md) для технічних деталей

**🇺🇦 Створено з ❤️ для українських контент-мейкерів**