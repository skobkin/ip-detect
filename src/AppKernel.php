<?php

namespace App;

use App\Controller\DetectorController;
use Symfony\Component\HttpFoundation\{Request, Response};
use Symfony\Component\HttpFoundation\Exception\RequestExceptionInterface;

final class AppKernel
{
    public const KERNEL_ROOT = __DIR__;
    public const PROJECT_ROOT = __DIR__.'/..';

    public static function handle(Request $request): void
    {
        $response = new Response();

        try {
            $response = static::handleRaw($request);
        } catch (RequestExceptionInterface $e) {
            $response->setStatusCode($e->getCode());
            $response->setContent($e->getMessage());
        } catch (\Throwable $e) {
            $response->setStatusCode(Response::HTTP_INTERNAL_SERVER_ERROR);
            $response->setContent('Server error');
        }

        $response->send();
    }

    public static function handleRaw(Request $request): Response
    {
        // It's a kind of routing you know...
        if ('/json' === $request->getPathInfo()) {
            $format = DetectorController::FORMAT_JSON;
        } elseif ('/plain' === $request->getPathInfo()) {
            $format = DetectorController::FORMAT_PLAIN;
        } elseif ('/' === $request->getPathInfo()) {
            $format = DetectorController::FORMAT_HTML;
        } else {
            return new Response('Resource not found', Response::HTTP_NOT_FOUND);
        }

        return DetectorController::detect($request, $format);
    }
}