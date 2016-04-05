<?php
try {
    convert($argc, $argv);
} catch (Exception $e) {
    echo $e->getMessage()."\n";
}


function error($msg)
{
    throw new Exception($msg);
}

function convert($argc, $argv)
{
    if ($argc < 2) {
        error('Usage '.$argv[0].' rconf path');
    }

    $rconfPath = $argv[1];
    if (!file_exists($rconfPath)) {
        error("$rconfPath not exists");
    }

    require($rconfPath);

    convertEnv($ENV);
    convertPrj($PRJ);
    convertSys($SYS);
}

function convertEnv($ENV)
{
    $data = $ENV['dev'];
    $data['NGX_EXEC_PREFIX'] = array(
        'ligang' => '/usr/local/bin/dexec nginx',
        'default' => 'sudo /usr/local/nginx/sbin/nginx',
    );

    save('var', $data);
}

function convertPrj($PRJ)
{
    $data = array();

    foreach ($PRJ['init_path'] as $key => $item) {
        $data['mkdir'][] = array(
            'dir' => $item['path'],
            'mask' => $item['mask'],
            'sudo' => $item['sudo'],
        );
    }

    if (isset($PRJ['webserver'])) {
        if ($PRJ['webserver'] == 'nginx') {
            $data['exec'][] = '${NGX_EXEC_PREFIX} -s reload';
        }
    }

    if (isset($PRJ['script'])) {
        foreach ($PRJ['script'] as $item) {
            $data['exec'][] = $item['path'];
        }
    }

    save('action', $data);
}

function convertSys($SYS)
{
    $data = array();

    foreach ($SYS as $skey => $sdata) {
        foreach ($sdata as $k => $item) {
            $key = $skey.'_'.$k;
            $data[$key] = $item;
        }
    }

    save('tpl', $data);
}

function save($key, $data)
{
    $path = dirname($_SERVER['argv'][1]).'/'.$key.'.json';
    $contents = json_encode($data, JSON_PRETTY_PRINT);
    $contents = str_replace('\\/', '/', $contents);

    file_put_contents($path, $contents);
}
