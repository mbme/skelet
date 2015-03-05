package storage

var rawData = []struct {
	Name       string
	Data       string
	Categories []Category
}{
	{"TEST",
		`# test

and so on

- so we can make it coool!
- and add other stuff
- like http links: http://google.com

\"OK GO\" we can say

this is so ~~bad~~ cool

First Header  | Second Header
------------- | -------------
Content Cell  | Content Cell
Content Cell  | Content Cell`,
		[]Category{"test", "tables"}},

	{"Advertisement :)",
		`---
__Advertisement :)__

- __[pica](https://nodeca.github.io/pica/demo/)__ - high quality and fast image
  resize in browser.
- __[babelfish](https://github.com/nodeca/babelfish/)__ - developer friendly
  i18n with plurals support and easy syntax.

You will like those projects!`,
		[]Category{"test", "images"}},

	{"Basic",
		`
# h1 Heading
## h2 Heading
### h3 Heading
#### h4 Heading
##### h5 Heading
###### h6 Heading


## Horizontal Rules

___

---

***`,
		[]Category{"test", "headings"}},

	{"Typographic replacements",
		`Enable typographer option to see result.

(c) (C) (r) (R) (tm) (TM) (p) (P) +-

test.. test... test..... test?..... test!....

!!!!!! ???? ,,

Remarkable -- awesome

\"Smartypants, double quotes\"

'Smartypants, single quotes'`,
		[]Category{"test", "typography"}},

	{"Emphasis",
		`
**This is bold text**

__This is bold text__

*This is italic text*

_This is italic text_

~~Deleted text~~

Superscript: 19^th^

Subscript: H~2~O

++Inserted text++

==Marked text==`,
		[]Category{"test", "typography"}},

	{"Blockquotes",
		`

> Blockquotes can also be nested...
>> ...by using additional greater-than signs right next to each other...
> > > ...or with spaces between arrows.
`,
		[]Category{"test", "typography"}},

	{"Lists",
		"Unordered\n" +
			"\n" +
			" + Create a list by starting a line with `+`, `-`, or `*`\n" +
			"+ Sub-lists are made by indenting 2 spaces:\n" +
			`- Marker character change forces new list start:
    * Ac tristique libero volutpat at
    + Facilisis in pretium nisl aliquet
    - Nulla volutpat aliquam velit
+ Very easy!

Ordered

1. Lorem ipsum dolor sit amet
2. Consectetur adipiscing elit
3. Integer molestie lorem at massa


1. You can use sequential numbers...` +
			"\n 1. ...or keep all the numbers as `1.`\n\n" +

			`Start numbering with offset:

57. foo
1. bar

`,
		[]Category{"test", "typography"}},

	{"Code",
		`
Inline ` + "`code`" + `

Indented code

    // Some comments
    line 1 of code
    line 2 of code
    line 3 of code


Block code \"fences\"

` + "```" + `
Sample text here...
` + "```" + `

Syntax highlighting

` + "``` js" + `
var foo = function (bar) {
  return bar++;
};

console.log(foo(5));
` + "```",
		[]Category{"test", "typography"}},

	{"Tables",
		`
| Option | Description |
| ------ | ----------- |
| data   | path to data files to supply the data that will be passed into templates. |
| engine | engine to be used for processing templates. Handlebars is the default. |
| ext    | extension to be used for dest files. |

Right aligned columns

| Option | Description |
| ------:| -----------:|
| data   | path to data files to supply the data that will be passed into templates. |
| engine | engine to be used for processing templates. Handlebars is the default. |
| ext    | extension to be used for dest files. |
`,
		[]Category{"test", "tables"}},

	{"Links",
		`
[link text](http://dev.nodeca.com)

[link with title](http://nodeca.github.io/pica/demo/ \"title text!\")

Autoconverted link https://github.com/nodeca/pica (enable linkify to see)
`,
		[]Category{"test", "typography"}},
	{"1 Большой текст",
		`
Компания Mozilla в своём блоге объявила о новых правилах безопасности, касающихся дополнений к их браузеру Firefox. В данный момент дополнение можно установить как из их коллекции AMO, так и с любого другого сайта. Из-за этой свободы для разработчиков дополнений участились случаи появления злонамеренных расширений. Они подменяют части страницы, вставляют на страницы рекламу или совершают другие неприятные действия. Поэтому в скором времени дополнения, которые можно установить в браузер, будут проходить процедуру обязательной проверки и подписи.

Дополнения, находящиеся в коллекции AMO, в любом случае проходят проверку и будут подписываться автоматически. Другие дополнения можно будет подавать на автоматическую проверку и последующую подпись, а если они не пойдут автоматическую проверку, будет возможность запросить проверку вручную. Переходный период займёт 12 недель. После этого на Release и Beta версии браузера нельзя будет устанавливать неподписанные дополнения.

Разработчики дополнений смогут тестировать неподписанные версии на ночных сборках или специальных редакциях браузера для разработчиков. Для тех, кто использует проприетарные непубличные дополнения для своих проектов, будет существовать возможность делать это и далее. Таким образом в Mozilla надеются повысить безопасность их браузера для пользователей. Эти планы пока не распространяются на другие проекты компании, Thunderbird или SeaMonkey.
`,
		[]Category{"test", "text"}},
	{"A Large Text",
		`You might not think Canada’s digital spies are on par with those in the US and UK—but rest assured, America’s northern neighbour is just as capable of perpetuating mass surveillance on a global scale. Case in point: at over 200 locations around the world, spies from Canada's cyberintelligence agency have been monitoring huge volumes of global internet traffic travelling across the internet's core.

​From these locations, Communications Security Establishment (CSE) can track who is accessing websites and files of interest. Its analysts can also log email addresses, phone numbers and even the content of unencrypted communications—and retain encrypted communication for later study, too—as well as intercept passwords and login details for later access to remote servers and websites.

​But perhaps more importantly, tapping into global internet traffic is a means for CSE to monitor, and also exploit, an ever growing list of digital threats, such as vulnerabilities in networks and computers and the spread of malware as well as botnets and the computers under their control. In the process, analysts can keep tabs on both friendly and foreign governments conducting covert cyber attacks and infiltration of their own.
`,
		[]Category{"test", "text"}},
}
