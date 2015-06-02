<form action="{{ urlFor "Upload.Img" ":id" 3 }}"
      class="dropzone"
      id="my-awesome-dropzone"
      enctype="multipart/form-data">
  	{{.xsrfdata}}
	<input type="file" name="file" multiple/>
</form>

