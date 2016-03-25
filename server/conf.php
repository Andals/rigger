<?php
/**
* @file conf_parser.php
* @brief parse rigger conf
* @author ligang
* @version 1.0
* @date 2014-11-19
 */

class Conf
{/*{{{*/
    const DEF_ITEM_KEY = 'default';

    const PARSE_FUNC_ARG  = 'ARG';
    const PARSE_FUNC_MATH = 'MATH';

    private $_env_conf = array();
    private $_prj_conf = array();
    private $_sys_conf = array();

    private $_user       = '';
    private $_run_params = array();

    public function __construct($params=array())
    {/*{{{*/
        $this->_user = $this->_getenv('USER');
    }/*}}}*/

    public function parseConf($params)
    {/*{{{*/
        $this->_run_params = $params;

        $conf_path = $this->getRunValue('rigger_conf');
        if(!file_exists($conf_path))
        {
            throw new Exception("$conf_path 不存在");
        }
        require($conf_path);

        $env = $this->getRunValue('env');
        if('' == $env || !isset($ENV[$env]))
        {
            throw new Exception("invalid env");
        }

        $this->_parseEnvConf($ENV, $env);
        $this->_parsePrjConf($PRJ);
        $this->_parseSysConf($SYS);
    }/*}}}*/
    public function parseTplContents($contents)
    {/*{{{*/
        return $this->_parseItem($contents);
    }/*}}}*/

    public function getRunValue($key)
    {/*{{{*/
        return isset($this->_run_params[$key]) ? $this->_run_params[$key] : null;
    }/*}}}*/
    public function getSysConf()
    {/*{{{*/
        return $this->_sys_conf;
    }/*}}}*/
    public function getSysConfItem($key)
    {/*{{{*/
        return isset($this->_sys_conf[$key]) ? $this->_sys_conf[$key] : null;
    }/*}}}*/
    public function getInitPath()
    {/*{{{*/
        return isset($this->_prj_conf['init_path']) ? $this->_prj_conf['init_path'] : null;
    }/*}}}*/
    public function getWebserver()
    {/*{{{*/
        return isset($this->_prj_conf['webserver']) ? $this->_prj_conf['webserver'] : null;
    }/*}}}*/
    public function getRunScript($run_point, $run_time)
    {/*{{{*/
        if(!isset($this->_prj_conf['script']))
        {
            return array();
        }

        $result = array();
        foreach($this->_prj_conf['script'] as $key => $item)
        {
            if(($item['run_point'] == $run_point) && ($item['run_time'] == $run_time))
            {
                $result[$key] = $item;
            }
        }
        return $result;
    }/*}}}*/

    private function _getenv($name)
    {/*{{{*/
        $result = Tool::runShellCmd("echo $$name", true);
        return $result['last_line'];
    }/*}}}*/

    private function _parseEnvConf($ENV, $env_key)
    {/*{{{*/
        if(!isset($ENV[$env_key]))
        {
            throw new Exception("env $env_key 不存在");
        }

        foreach($ENV[$env_key] as $key => $item)
        {
            $this->_env_conf[$key] = $this->_parseItem($item);
        }
    }/*}}}*/
    private function _parsePrjConf($PRJ)
    {/*{{{*/
        $this->_prj_conf['webserver'] = $this->_parseItem($PRJ['webserver']);
        if(isset($PRJ['init_path']))
        {
            foreach($PRJ['init_path'] as $key => $data)
            {
                foreach($data as $k => $v)
                {
                    $this->_prj_conf['init_path'][$key][$k] = $this->_parseItem($v);
                }
            }
        }

        if(isset($PRJ['script']))
        {
            foreach($PRJ['script'] as $key => $data)
            {
                foreach($data as $k => $v)
                {
                    $this->_prj_conf['script'][$key][$k] = $this->_parseItem($v);
                }
            }
        }
    }/*}}}*/
    private function _parseSysConf($SYS)
    {/*{{{*/
        foreach($SYS as $sys_key => $sys_item)
        {
            foreach($sys_item as $key => $item)
            {
                foreach($item as $k => $v)
                {
                    $this->_sys_conf[$sys_key][$key][$k] = $this->_parseItem($v);
                }
            }
        }
    }/*}}}*/
    private function _parseReplaceUseFunc($func_name, $replace_name)
    {/*{{{*/
        $result = '';

        switch($func_name)
        {
        case self::PARSE_FUNC_ARG:
            if(!isset($this->_run_params[$replace_name]))
            {
                throw new Exception("arg $replace_name 不存在");
            }
            $result = $this->_run_params[$replace_name];
            break;
        case self::PARSE_FUNC_MATH:
            $result = $this->_math($replace_name);
            break;
        }

        return $result;
    }/*}}}*/
    private function _parseItem($item)
    {/*{{{*/
        $line = '';
        if(is_array($item))
        {
            $line = array_key_exists($this->_user, $item) ? $item[$this->_user] : $item[self::DEF_ITEM_KEY];
        }
        else
        {
            $line = $item;
        }

        if(preg_match_all('/__([A-Z]+)__\(([^)]+)\)/', $line, $matches))
        {
            list($search, $func, $replace) = $matches;
            foreach($replace as $index => $name)
            {
                $replace[$index] = $this->_parseReplaceUseFunc($func[$index], $name);
            }
            $line = str_replace($search, $replace, $line);
        }
        if(preg_match_all('/\${([^}]+)}/', $line, $matches))
        {
            list($search, $replace) = $matches;
            foreach($replace as $index => $name)
            {
                $replace[$index] = isset($this->_env_conf[$name]) ? $this->_env_conf[$name] : $this->_getenv($name);
            }
            $line = str_replace($search, $replace, $line);
        }

        return $line;
    }/*}}}*/

    private function _math($expression)
    {/*{{{*/
        $result = '';
        if(preg_match('#([^+\-*/]+)([+\-*/])([^+\-*/]+)#', $expression, $matches))
        {
            $v1 = $this->_parseItem($matches[1]);
            $v2 = $this->_parseItem($matches[3]);
            $op = $matches[2];
            switch($op)
            {
                case '+':
                    $result = $v1 + $v2;
                    break;
                case '-':
                    $result = $v1 - $v2;
                    break;
                case '*':
                    $result = $v1 * $v2;
                    break;
                case '/':
                    $result = $v1 / $v2;
                    break;
            }
        }
        return $result;
    }/*}}}*/
}/*}}}*/
