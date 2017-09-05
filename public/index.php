<?php

namespace App;

use \Symfony\Component\HttpFoundation\Request;

$loader = require __DIR__.'/../vendor/autoload.php';

$request = Request::createFromGlobals();

AppKernel::handle($request);