<?php
$ENV = array
(/*{{{*/
    'dev' => array
    (/*{{{*/
        'IS_DEV'    => 'true',
        'USER'      => '${USER}',
        'PRJ_NAME'  => 'demo',
        'PRJ_HOME'  => '__ARG__(prj_home)',

        'FRONT_DOMAIN'        => '${USER}.demo.com',
        'FRONT_ACCESS_LOG'    => '${PRJ_NAME}_front_${USER}.log',
        'FRONT_ERROR_LOG'     => '${PRJ_NAME}_front_${USER}_error.log',
        'FRONT_HTTP_CONF_TPL' => '${PRJ_HOME}/rigger/client/tpl/tpl_front_httpd.conf.ngx',
        'FRONT_HTTP_CONF_DST' => '${PRJ_HOME}/config/http/${USER}_front_http.conf.ngx',
        'FRONT_HTTP_CONF_LN'  => '/usr/local/nginx/conf/include/${FRONT_DOMAIN}.conf',

        'ACCESS_LOG_BUFFER' => '1',
        'NGX_PORT'          => '80',
        'GO_PORT'           => '__MATH__(7000+${UID})',

        'SERVER_CONF_TPL' => '${PRJ_HOME}/rigger/client/tpl/tpl_server_conf.php',
        'SERVER_CONF_DST' => '${PRJ_HOME}/config/server/${USER}_server_conf.php',
        'SERVER_CONF_LN'  => '${PRJ_HOME}/config/server_conf.php',
    ),/*}}}*/
);/*}}}*/

$PRJ = array
(/*{{{*/
    'webserver' => 'nginx',
    'init_path' => array
    (/*{{{*/
        'tmp' => array(
            'path' => '${PRJ_HOME}/tmp/',
            'mask' => '777',
            'sudo' => false
        ),
        'logs/front' => array(
            'path' => '${LOG_PATH}/front',
            'mask' => '777',
            'sudo' => false
        ),
        'logs/task' => array(
            'path' => '${LOG_PATH}/task',
            'mask' => '755',
            'sudo' => false
        ),
    ),/*}}}*/
);/*}}}*/

$SYS = array
(/*{{{*/
    'front' => array
    (/*{{{*/
        'server_conf' => array(
            'tpl'  => '${SERVER_CONF_TPL}',
            'dst'  => '${SERVER_CONF_DST}',
            'ln'   => '${SERVER_CONF_LN}',
            'sudo' => false,
        ),
        'front_http_conf' => array(
            'tpl'  => '${FRONT_HTTP_CONF_TPL}',
            'dst'  => '${FRONT_HTTP_CONF_DST}',
            'ln'   => '${FRONT_HTTP_CONF_LN}',
            'sudo' => true,
        ),
    ),/*}}}*/
);/*}}}*/
