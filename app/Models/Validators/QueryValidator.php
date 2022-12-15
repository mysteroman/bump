<?php namespace Models\Validators;

use Zephyrus\Application\Rule;

class QueryValidator
{
    public static function validate($form): bool
    {
        $form->field('placeId')->validate([
            Rule::notEmpty(),
            CustomRule::placeId(),
            Rule::maxLength(255)
        ], true);
        $form->field('route')->validate([
            Rule::notEmpty(),
            Rule::maxLength(255)
        ]);
        if (!$form->verify()) return false;
        $obj = $form->buildObject();
        return !empty($obj->placeId) || !empty($obj->route) && !(!empty($obj->placeId) && !empty($obj->route));
    }
}