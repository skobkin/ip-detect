# Minimalistic IP address detector app

[![Codeship Status for skobkin/ip-detect](https://app.codeship.com/projects/b2375af0-c0e5-0136-9e74-26b840c766cc/status?branch=master)](https://app.codeship.com/projects/313611)

[Try it](https://ip.skobk.in/) ([JSON](https://ip.skobk.in/json), [Plaintext](https://ip.skobk.in/plain))

## Installation

```bash
# For usage (from Packagist)
composer create-project skobkin/ip-detect
# For development (from Git)
git clone git@bitbucket.org:skobkin/ip-detect.git
cd ip-detect && composer install
```


## Running

```bash
# From the project root directory
php -S localhost:8000 -t public public/index.php
```
