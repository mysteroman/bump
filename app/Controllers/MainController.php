<?php namespace Controllers;


use Models\Brokers\DataBroker;
use Models\Services\DataService;
use Models\Validators\QueryValidator;

class MainController extends Controller
{
    public function initializeRoutes()
    {
        $this->get('/', 'index');
        $this->get('/ranking', 'ranking');
        $this->get('/api', 'api');
    }

    public function index()
    {
        return $this->render('index');
    }

    public function ranking()
    {
        return $this->render('ranking', DataService::getRankings());
    }

    public function api()
    {
        $form = $this->buildForm();
        return $this->json((object)[
            'route' => DataService::find($form)
        ]);
    }
}