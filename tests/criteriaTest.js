/* Criteria add test
*  opens criteria list
*   then opens criteria add list
*   then adds criterion
*   then returns to criteria list
*   then checks the list
*/

phantom.injectJs("header.js"); // opens header
 


// page.thenOpen(base_url+'/adm/crit/list', function(){
//     this.echo(this.getCurrentUrl()+' opened  ', 'PARAMETER');
//     this.waitForUrl(base_url+"/adm/crit/list", function then() {
//         this.echo(this.getTitle() + ' | sucessfully opened Criteria','INFO');
//     }, function timeout() {
//         this.echo(this.getCurrentUrl());
//         this.capture('screen_logs/critOpen.png');
//         this.die(this.getTitle() + '| Criteria list Opening  failed ','ERROR');
//     });
// });

// page.then(function(){
//     this.click('a[href="/adm/crit/add"]');
//     this.echo(this.getCurrentUrl()+' opened','PARAMETER');
//     this.waitForUrl(base_url+'/adm/crit/add', function then(){
//    //     this.echo("Adding a criterion" , 'INFO');
//     }, function timeout() {
//         this.echo(this.getCurrentUrl(),'ERROR');
//         this.capture('screen_logs/critAddError.png');
//         this.die(this.getTitle() + '| Criterion adding  link open error','ERROR');
//     });
 
// });

// page.then(function(){
//       this.fill('form[action="/adm/crit/add"]', {
//         "title": 'testcsdf',
//         "dep":['2'],

//     },true);
//       this.click('button[type="submit"]');
       
// });


// page.then(function(){
// 	this.echo('Adding a criterion');
// });

// page.then(function(){
	
// 	this.waitForUrl(base_url+"/adm/crit/list", function then() {
	
//         this.echo(this.getTitle() + ' | Added a new criterion','INFO');
       
//     }, function timeout() {
//         this.echo(this.getCurrentUrl());
//         this.capture('screen_logs/critOpenBack.png');
//         this.die(this.getTitle() + '| Criterion adding failed ','ERROR');
//     });


// });

page.thenOpen(base_url+"/adm/crit/list" , function(){
 this.mouse.click("#headingOne");
        
     this.waitForSelector('a[href^="/adm/crit/"]', function(){
            this.wait(2000);
        this.capture('screen_logs/whh.png')
           
         

   
 });
// this.click('a[href$="/remove"]');
  // 

});

page.then(function(){
  this.click('.btn-sm.btn-success')
    this.capture('screen_logs/wh.png');
});

page.waitForUrl(base_url+"/adm/crit/", function(){
    this.echo(this.getCurrentUrl());    
    this.capture('screen_logs/wh.png');
    this.captureSelector('screen_logs/selector.png', 'form[action^="/adm/crit"]');

});

page.then(function(){
  this.capture('screen_logs/mid.png')
});
page.then(function(){
    this.echo(this.getCurrentUrl());
  });

page.then(function(){

    this.fill('form[action^="/adm/crit"]', {
        'title': 'I am watching you',
        'dep':['2'],
        },true);
    // this.mouse.click('.btn-primary');
    this.capture('screen_logs/after.png')
     this.waitForSelector('.panel.panel-info', function(){
            this.capture('screen_logs/addCritt.png');
      });

});

            // this.wait(2000);
    


page.then(function(){
   
    
  
  
  // this.mouse.doubleclick('.btn-primary');
  this.wait(500);
    this.capture('screen_logs/wee.png')
});
// page.wait(2000);
page.then(function(){


});

page.then(function(){
    this.echo(this.getCurrentUrl());
    this.capture('screen_logs/weee.png')
  });

page.run();