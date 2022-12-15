<?php namespace Models\Validators;

use Zephyrus\Application\Rule;

class QueryValidator
{
    public static function validate($form): bool
    {
        $form->field('route')->validate([
            Rule::notEmpty(),
            Rule::maxLength(255)
        ]);
        $form->field('placeId')->validate([
            Rule::notEmpty(),
            CustomRule::placeId(),
            Rule::maxLength(255)
        ], true);
        
        return $form->verify();
    }
}