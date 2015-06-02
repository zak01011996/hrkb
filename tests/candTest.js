
phantom.injectJs("header.js"); 

page.thenOpen(base_url+'/adm/candidates/list', function(){ 
    this.echo(this.getCurrentUrl()+' opened  ', 'PARAMETER');
    this.waitForUrl(base_url+"/adm/candidates/list", function then() {
        this.echo(this.getTitle() + ' | sucessfully opened candidate','INFO');
    }, function timeout() {
        this.echo(this.getCurrentUrl());
        this.capture('screen_logs/candidateOpen.png');
        this.die(this.getTitle() + '| candidate list Opening  failed ','ERROR');
    });
});


page.then(function(){
    this.click('a[href="/adm/candidates/add"]');
    this.echo(this.getCurrentUrl()+' opened','PARAMETER');
    this.waitForUrl(base_url+'/adm/candidates/add', function then(){
        this.echo("Adding candidates" , 'INFO');
    }, function timeout() {
        // this.echo(this.getCurrentUrl(),'ERROR');
        this.capture('screen_logs/departAddError.png');
        this.die(this.getTitle() + '| candidate adding  link open error','ERROR');
    });
 
});

page.then(function(){
    // this.echo(this.getCurrentUrl(),'ERROR');
      this.fill('form[action="/adm/candidates/add"]', {
        "name": 'candidadtecsName',
        "lname": 'Lname?',
        "fname": 'Fname?',
        "phone":'1239013819830',
        "email": 'name@yourtld.com',
        "note":'Checking note area ',
        "addr":'Adress of someone will be written',
        "married": true,
        "depId": ['1'],
        "salary":12310,
        "currency": 'curr',
    },true);
      this.waitForSelector('.table.table-striped', function(){
            this.capture('screen_logs/adddCandd.png');
      });
       
});

page.then(function(){
    // this.echo(this.getCurrentUrl(),'ERROR');
    this.click('.btn-warning');
    this.echo(this.getCurrentUrl());

   this.waitForSelector('form[action="/adm/candidates/add"]', function() {
    this.fill('form[action="/adm/candidates/add"]',{
         "name": 'ChangedName',
        "lname": 'ChLname?',
        "fname": 'ChFname?',
        "phone":'12239013819830',
        "email": 'chname@yourtld.com',
        "note":'Checking note Ch ',
        "addr":'Checking Adress of someone will be written',
        "married": false,
        "depId": ['1'],
        "salary":122310,
        "currency": 'Checking',

    },true);


     this.waitForSelector('.table.table-striped', function(){
            this.capture('screen_logs/adddCandd.png');
      });

});


});


page.then(function(){
    this.click('img');
  
});

page.then(function(){
    this.echo(this.getCurrentUrl(),'ERROR');
   
    this.capture('screen_logs/ajskdhaksjdhasjd.png')
 this.fill('form[action$="/contacts/add"]', {
            'name': "candidadtecsName",
            'value':"somevalue"
        },false);

 // this.capture('screen_logs/asdaksjd.png')

    this.click('#add_cont');
   

}).wait(5000);


page.then(function(){
    this.click('#show_dropzone');

      this.waitForSelector('#my-awesome-dropzone',function(){
       
       this.click('#my-awesome-dropzone');
      this.capture('screen_logs/upload.png')

        this.echo('Give me good answer', 'ERROR')

    });
    
    this.fill('#rating_add_form',{
        'crit': ['3'],
        'value':['2']
      },true);
      this.click('button[type="submit"]');
      this.wait(2000);
   
        this.echo(this.getCurrentUrl());
    this.capture('screen_logs/contactAdded.png');
    this.fill('form[action^="/adm/comments/"]', {
        'text': "Test text for special comment"
 
    },true);
    this.echo(this.getCurrentUrl());
    this.capture('screen_logs/addCand.png')
    // this.click('b')
});


page.run();