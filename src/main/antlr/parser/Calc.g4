grammar Calc;

prog: stat+;

stat: expr ';'                      # cal_stat
    | ID '=' expr ';'               # assign
    | 'print' '(' expr ')' ';'      # print
    ;

expr: <assco=right>expr POW expr    # pow_expr
    | expr op=(MUL|DIV|MOD) expr    # md_expr
    | expr op=(ADD|SUB) expr        # as_expr
    | sign=(ADD|SUB)? NUM           # number
    | ID                            # id
    | '(' expr ')'                  # parens
    ;

NUM: INT
   | FLOAT
   ;

MUL: '*';
DIV: '/';
ADD: '+';
SUB: '-';
POW: '^';
MOD: '%';

ID: [a-zA-Z]+[0-9a-zA-Z]*;

ZERO: '0';

INT: [1-9][0-9]*
    | ZERO
    ;

FLOAT: INT '.' [0-9]+
    ;

COMMENT_LINE: '//' .*? ('\r'|'\n') -> skip;
COMMENT_BLOCK: '/*' .*? '*/'     -> skip;

WS: [ \t\r\n] -> skip;