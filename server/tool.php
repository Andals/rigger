<?php
class Tool
{/*{{{*/
    public static function runShellCmd($cmd, $fetch_output=false)
    {/*{{{*/
        $return_var = 0;
        $last_line  = '';
        $output     = '';

        ob_start();
        $last_line  = system($cmd, $return_var);
        $cmd_output = ob_get_clean();

        if($fetch_output)
        {
            $output = $cmd_output;
        }
        else
        {
            echo $cmd_output;
        }

        return array(
            'return_var' => $return_var,
            'last_line'  => $last_line,
            'output'     => $output,
        );
    }/*}}}*/
    public static function runShell($shell, $params=array(), $fetch_output=false)
    {/*{{{*/
        $params = implode(' ', $params);
        $cmd    = "$shell $params";
        return self::runShellCmd($cmd, $fetch_output);
    }/*}}}*/
}/*}}}*/
