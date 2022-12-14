<?php namespace Models\Services;


use Models\Brokers\DataBroker;
use Models\Validators\QueryValidator;

class DataService
{
    public static function getRankings(): array
    {
        $broker = new DataBroker();
        return [
            'best' => $broker->getBestRankings(),
            'worst' => $broker->getWorstRankings()
        ];
    }

    public static function find($form): ?\stdClass
    {
        if (!QueryValidator::validate($form)) return null;
        $obj = $form->buildObject();
        if (!empty($obj->placeId)) return (new DataBroker())->findByPlaceId($obj->placeId);
        return (new DataBroker())->findByName($obj->route);
    }
}