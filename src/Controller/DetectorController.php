<?php

namespace App\Controller;

use App\AppKernel;
use Slim\Views\PhpRenderer;
use Symfony\Component\HttpFoundation\{JsonResponse, Request, Response};

final class DetectorController
{
    public const FORMAT_JSON = 'json';
    public const FORMAT_HTML = 'html';
    public const FORMAT_PLAIN = 'plaintext';

    private const TEMPLATES_PATH = AppKernel::PROJECT_ROOT.'/templates';

    public static function detect(Request $request, string $format = 'html'): Response
    {
        $clientData = [
            'ip_address' => $request->getClientIp(),
            'locale' => $request->getLocale(),
            'preferred_language' => $request->getPreferredLanguage(),
        ];

        if (static::FORMAT_JSON === $format) {
            return new JsonResponse($clientData);
        } elseif (static::FORMAT_PLAIN === $format) {
            return new Response($request->getClientIp());
        }

        $templater = new PhpRenderer(static::TEMPLATES_PATH);
        $templater->setAttributes($clientData);

        return new Response($templater->fetch('/client_data.php', $clientData));
    }
}