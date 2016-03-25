<?php
/**
* @file rigger.php
* @brief rigger
* @author ligang
* @version 1.0
* @date 2014-11-19
 */

class Rigger
{/*{{{*/
    const ALL_SYS_KEY = 'all';

    const RUN_TIME_PRE  = 'pre';
    const RUN_TIME_POST = 'post';

    private $_cmd_list    = array();
    private $_rigger_conf = null;

    public function __construct()
    {/*{{{*/
        $this->_initCmdList();
    }/*}}}*/

    public function run($params=array())
    {/*{{{*/
        $this->_rigger_conf = new Conf();
        $this->_rigger_conf->parseConf($params);

        $cmd = $this->_rigger_conf->getRunValue('cmd');
        if(!isset($this->_cmd_list[$cmd]))
        {
            $this->_printCmdList();

            throw new Exception("cmd 不正确\n");
        }

        $conf = $this->_cmd_list[$cmd];
        $func = $conf['func'];

        $this->_preCmd($cmd);
        $this->$func();
        $this->_postCmd($cmd);
    }/*}}}*/


    private function _initCmdList()
    {/*{{{*/
        $this->_cmd_list = array(
            'conf' => array(
                'desc' => '生成实际配置',
                'func' => '_conf',
            ),
            'restart' => array(
                'desc' => '重启服务',
                'func' => '_restart',
            ),
        );
    }/*}}}*/
    private function _printCmdList()
    {/*{{{*/
        $result = "\nCmd list:\n";
        foreach($this->_cmd_list as $cmd => $item)
        {
            $result.= "$cmd\t".$item['desc']."\n";
        }

        echo $result;
    }/*}}}*/

    private function _preCmd($cmd)
    {/*{{{*/
        $script = $this->_rigger_conf->getRunScript($cmd, self::RUN_TIME_PRE);
        if(!empty($script))
        {
            $this->_runScript($script);
        }
    }/*}}}*/
    private function _postCmd($cmd)
    {/*{{{*/
        $script = $this->_rigger_conf->getRunScript($cmd, self::RUN_TIME_POST);
        if(!empty($script))
        {
            $this->_runScript($script);
        }
    }/*}}}*/

    private function _conf()
    {/*{{{*/
        $sys_key = $this->_rigger_conf->getRunValue('sys');
        if(self::ALL_SYS_KEY == $sys_key)
        {
            foreach($this->_rigger_conf->getSysConf() as $sys_item)
            {
                $this->_genConfFile($sys_item);
            }
        }
        else
        {
            $sys_item = $this->_rigger_conf->getSysConfItem($sys_key);
            if(is_null($sys_item))
            {
                throw new Exception("sys $sys_key not exists");
            }

            $this->_genConfFile($sys_item);
        }

        $init_path = $this->_rigger_conf->getInitPath();
        if(!is_null($init_path))
        {
            foreach($init_path as $item)
            {
                $this->_initSinglePath($item);
            }
        }
    }/*}}}*/
    private function _restart()
    {/*{{{*/
        $webserver = $this->_rigger_conf->getWebserver();
        if('' != $webserver)
        {
            $this->_restartWebServer($webserver);
        }
    }/*}}}*/


    private function _genConfFile($sys_item)
    {/*{{{*/
        foreach($sys_item as $item)
        {
            $this->_genSingleConfFile($item);
        }
    }/*}}}*/
    private function _genSingleConfFile($item)
    {/*{{{*/
        $tpl = trim($item['tpl']);
        if(!file_exists($tpl))
        {
            throw new Exception("tpl $tpl not exists");
        }

        $tpl_contents = file_get_contents($tpl);
        $dst_contents = $this->_rigger_conf->parseTplContents($tpl_contents);

        $dst     = trim($item['dst']);
        $dst_dir = dirname($dst);
        if(!is_dir($dst_dir))
        {
            throw new Exception("dst $dst_dir not exists");
        }
        file_put_contents($dst, $dst_contents);

        $ln     = trim($item['ln']);
        $ln_dir = dirname($ln);
        if(!is_dir($ln_dir))
        {
            throw new Exception("ln $ln_dir not exists");
        }

        $sudo = trim($item['sudo']);
        if(file_exists($ln))
        {
            $cmd = "";
            if($sudo)
            {
                $cmd = "sudo ";
            }
            $cmd.= "rm -f $ln";
            Tool::runShellCmd($cmd);
        }

        $cmd = "";
        if($sudo)
        {
            $cmd.= "sudo ";
        }
        $cmd.= "ln -s $dst $ln";
        Tool::runShellCmd($cmd);
    }/*}}}*/

    private function _initSinglePath($item)
    {/*{{{*/
        $sudo = trim($item['sudo']);
        $path = trim($item['path']);
        if(!is_dir($path))
        {
            $cmd = "";
            if($sudo)
            {
                $cmd.= "sudo ";
            }
            $cmd.= "mkdir -p $path";
            Tool::runShellCmd($cmd);
        }

        $mask = trim($item['mask']);
        $cmd  = "";
        if($sudo)
        {
            $cmd.= "sudo ";
        }
        $cmd.= "chmod -R $mask $path";
        Tool::runShellCmd($cmd);
    }/*}}}*/

    private function _restartWebServer($webserver)
    {/*{{{*/
        switch($webserver)
        {
        case 'nginx':
            $this->_restartNginx();
            break;
        case 'apache':
            $this->_restartApache();
            break;
        default:
            throw new Exception("invalid webserver $webserver");
        }
    }/*}}}*/
    private function _restartNginx()
    {/*{{{*/
        $cmd    = 'sudo /usr/local/nginx/sbin/nginx -t';
        $result = Tool::runShellCmd($cmd);
        if(0 != $result['return_var'])
        {
            throw new Exception("nginx conf error");
        }

        $cmd = 'sudo /usr/local/nginx/sbin/nginx -s reload';
        Tool::runShellCmd($cmd);

        echo "nginx restart successful\n"; 
    }/*}}}*/
    private function _restartApache()
    {/*{{{*/
        $cmd = 'sudo /usr/local/apache2/bin/apachectl -t';
        $result = Tool::runShellCmd($cmd);
        if(0 != $result['return_var'])
        {
            throw new Exception("apache conf error");
        }

        $cmd = 'sudo /usr/local/apache2/bin/apachectl graceful';
        Tool::runShellCmd($cmd);

        echo "apache restart successful\n"; 
    }/*}}}*/

    private function _runScript($script)
    {/*{{{*/
        foreach($script as $key => $item)
        {
            $this->_runSingleScript($item);
        }
    }/*}}}*/
    private function _runSingleScript($item)
    {/*{{{*/
        $path = trim($item['path']);
        if(file_exists($path))
        {
            Tool::runShell($path);
        }
    }/*}}}*/
}/*}}}*/
