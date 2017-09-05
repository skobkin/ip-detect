<?php /** @var \Slim\Views\PhpRenderer $this */ ?>
<!DOCTYPE html>
<html>
    <head>
        <title>IP address</title>
    </head>

    <body>
        <h1><?php echo htmlspecialchars($this->getAttribute('ip_address')); ?></h1>

        <pre>
locale: <?php echo htmlspecialchars($this->getAttribute('locale')); ?>

preferred language: <?php echo htmlspecialchars($this->getAttribute('preferred_language')); ?>
        </pre>
    </body>
</html>