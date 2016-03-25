<?php
require_once(__DIR__.'/rigger.php');
require_once(__DIR__.'/conf.php');
require_once(__DIR__.'/tool.php');

$params = array();
$args   = array_slice($argv, 1);
foreach($args as $arg)
{
    $arg   = explode('=', $arg);
    $key   = trim($arg[0]);
    $value = trim($arg[1]);
    $params[$key] = $value;
}

$rigger = new Rigger();
try
{
    $rigger->run($params);
}
catch(Exception $e)
{
    echo $e->getMessage()."\n";
}
