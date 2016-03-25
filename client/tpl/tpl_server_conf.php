<?php
namespace Demo\Config;

class ServerConf extends \Demo\Config\Lib\Dev
{/*{{{*/
    public function getPrjHome()
    {/*{{{*/
        return '${PRJ_HOME}';
    }/*}}}*/
}/*}}}*/
