/*      Add user Test:
*       Adds user
*       checks with existing user
*       
*/

phantom.injectJs("header.js"); //import header

/* opens users list
*  prints current status
*  takes screnshot on error > userListOpen.png
*    
*/

page.thenOpen(base_url+'/adm/users/list', function(){
    this.echo(this.getCurrentUrl()+' opened  ', 'PARAMETER');
    this.waitForUrl(base_url+"/adm/users/list", function then() {
        this.echo(this.getTitle() + ' | sucessfully opened users list','INFO');
    }, function timeout() {
        this.echo(this.getCurrentUrl());
        this.capture('screen_logs/userListOpen.png');
        this.die(this.getTitle() + '| user list Opening  failed ','ERROR');
    });
});


/* Adds user
*  takes screenshot on error addUserError.png
*/

page.then(function(){

    this.click('a[href="/adm/users/add"]');
    this.echo(this.getCurrentUrl()+' opened','PARAMETER');
    this.waitForUrl(base_url+'/adm/users/add', function then(){
        this.echo("Adding user" , 'INFO');
    }, function timeout() {
        this.echo(this.getCurrentUrl(),'ERROR');
        this.capture('screen_logs/addUserError.png');
        this.die(this.getTitle() + '| user adding error','ERROR');
    });
 
});

// fills  form


page.then(function(){
      this.fill('form[action="/adm/users/add"]', {
        "login": 'SampleUserNameLogin',
        "name":'SampleNameForTest',
        "pass": 'password1',
        "passc": 'password1',
        "role":['1'],   

    },true);
       this.echo('filled form', 'INFO');
});

// prints error when error selector exist
// screenshot on error >  errorError.png
page.then(function(){
    if (this.exists('.control-label')) {
        this.echo('ERROR ! '+ this.getHTML('.control-label') , 'ERROR');
        this.echo(this.getElementAttribute('label.control-label'));
        // this.capture('screen_logs/screensho222t.png')

    }
    else {
        this.capture('screen_logs/errorError.png')
        this.echo('shouldnt be any error', 'INFO');
    }


});

// loads page again
page.then(function(){

    this.reload(function() {
        this.echo("loaded again");
        
    });
})

//adds existing user

page.then(function(){
this.capture('screen_logs/reloadpage.png')
    this.echo('checking with existing user ', 'INFO');
      this.fill('form[action="/adm/users/add"]', {
        "login": 'admin',
        "name":'userName1',
        "pass": 'password1',
        "passc": 'password11',
        "role":['1'],   

    },true);
       this.echo('filled form', 'INFO');
          //this.click('button[type="submit"]');


});

// screenshot on error addExistUserEr.png
page.then(function(){
    if (this.exists('label.control-label')) {
        this.echo('ERROR ! '+ this.getHTML('label.control-label') , 'ERROR');
        this.echo(this.getElementAttribute('label.control-label'));
        this.capture('screen_logs/addExistUserEr.png')
    }
    else {
        this.capture('screen_logs/sasdacreensho222t.png')
        this.echo('shouldnt be any error', 'INFO');
    }


});



page.then(function(){
    this.echo('Terminating . . . ', 'PARAMETER');
    // this.capture('screen_logs/screenshooot.png')
});

page.run();
