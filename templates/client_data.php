<?php /** @var \Slim\Views\PhpRenderer $this */ ?>
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>IP address</title>
    </head>

    <body>
        <h1><?php echo htmlspecialchars($this->getAttribute('ip_address')); ?></h1>

        <pre>
locale: <?php echo htmlspecialchars($this->getAttribute('locale')); ?>

preferred language: <?php echo htmlspecialchars($this->getAttribute('preferred_language')); ?>
        </pre>

        <br />

        <a href="https://skobk.in/">skobk.in</a>
    </body>
</html>