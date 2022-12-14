<?php namespace Models\Validators;

use Zephyrus\Application\Rule;

class CustomRule
{
    public static function placeId(): Rule
    {
        return new Rule(function($data) {
            return preg_match('/^[a-zA-Z0-9_\-]+$/', $data);
        });
    }
}
