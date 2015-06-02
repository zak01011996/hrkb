/* Department Test 
* Tests Department section:
*       opens departments
*       adds department
*       returns to department list
*       then adds another department
*       exits
*       input : 
*           department name1
*           department name2
*       output :
*           added two department
*       
*/

phantom.injectJs("header.js"); // import header

/*  Opens departments list
*   prints sucess or error message
* when error takes screenshot    > >depOpen.png
*/

page.thenOpen(base_url+'/adm/dep/list', function(){ 
    this.echo(this.getCurrentUrl()+' opened  ', 'PARAMETER');
    this.waitForUrl(base_url+"/adm/dep/list", function then() {
        this.echo(this.getTitle() + ' | sucessfully opened Departments','INFO');
    }, function timeout() {
        this.echo(this.getCurrentUrl());
        this.capture('screen_logs/depOpen.png');
        this.die(this.getTitle() + '| Departments list Opening  failed ','ERROR');
    });
});

/*  Opens adding department
*   prints current status on console
*   takes screenshot >  departAddError.png
*/

page.then(function(){
    this.click('a[href="/adm/dep/add"]');
    this.echo(this.getCurrentUrl()+' opened','PARAMETER');
    this.waitForUrl(base_url+'/adm/dep/add', function then(){
        this.echo("Adding department" , 'INFO');
    }, function timeout() {
        this.echo(this.getCurrentUrl(),'ERROR');
        this.capture('screen_logs/departAddError.png');
        this.die(this.getTitle() + '| Department adding  link open error','ERROR');
    });
 
});

/* Adds new department
*/
page.then(function(){
      this.fill('form[action="/adm/dep/add"]', {
        "title": 'FirstDepartmentForTest',

    },true);
       
});
page.then(function(){
    this.echo('trying to add a Department', 'INFO');
     

});

/* returns to department list
*  prints current status
* exits*/

page.then(function(){
    
    this.waitForUrl(base_url+"/adm/dep/list", function then() {
    
        this.echo(this.getTitle() + ' | sucessfully returned to Departments list','INFO');
    }, function timeout() {
        this.echo(this.getCurrentUrl());
        this.capture('screen_logs/depOpenBack.png');
        this.die(this.getTitle() + '| Departments list return  failed ','ERROR');
    });


});



page.thenOpen(base_url+"/adm/dep/add", function(){
        this.fill('form[action="/adm/dep/add"]', {
                "title": 'NewSampledepsmentForTest',

            },true);
               this.echo('trying to add another Department', 'INFO');
});

page.then(function(){
    this.echo(this.getCurrentUrl()+' opened','PARAMETER');

});

page.then(function(){
    this.capture('screen_logs/whatisi.png');
     if (this.exists('.control-label')) {
        this.echo('ERROR ! '+ this.fetchText('.control-label') , 'ERROR');
        this.echo(this.getElementAttribute('label.control-label'));
        this.capture('screen_logs/DeparmentAddError.png')
    }
    else {
        this.capture('screen_logs/DeparmentAddError.png')
        this.echo('shouldnt be any error', 'INFO');
    }
});


page.run();
