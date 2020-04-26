%{
package goyacc
import(
    "fmt"
    "strings"
    "bufio"
)
const(
    Debug = 4
    ErrorVerbose = true
)
%}

%union {
    empty struct{}
    s string
    i interface{}
    m map[string]interface{}
}

%token <empty> EQ NE GT GE LT LE
%token <s> KEY
%token <i> VAL

%type<s> key
%type<i> val

%%
expr:
    key EQ val 
    {
    fmt.Println("eq", $1, $3)
    }
|   key NE val
    {
    fmt.Println("ne", $1, $3)
    }

key:
    KEY{
    $$ = $1
    }

val:
    VAL
    {
    $$ = $1
    }

%%

func Parse(input string){
	//
	strReader := strings.NewReader(strings.TrimSpace(input))
    //
    yyParse(&line{ input: input, buf: bufio.NewReader(strReader) })
}

