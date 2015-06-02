var page = require('casper').create();
var username="admin"
var password="1"
var base_url='http://localhost:8080'
// First function
page.start(base_url+'/login', function () {
    this.page.injectJs('https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js');
    this.viewport(1024,768);
    this.echo(this.getTitle(),'PARAMETER');
});


page.then(function(){
this.fill('form[action="/login"]', {
        "login": username,
        "password": password,
    },true);
	//this.click('button[type="submit"]');
});

page.then(function(){
	this.waitForUrl('http://localhost:8080/adm/candidates', function then() {
    this.echo('Login sucessfull with \nusername : ' + username + '\npassword: ' + password, 'PARAMETER');
    },function timeout() {
        this.capture('screen_logs/loginOpenError.png');
        this.echo(this.getHTML('control-label'));
        this.die('Login unsucessfull with \nusername : ' + username + '\npassword: ' + password,'ERROR');
    });
	//this.click('button[type="submit"]');
	this.capture('screen_logs/shu.png');
});


