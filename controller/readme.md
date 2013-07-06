Controller
==========

Controller is as rails-inspired implementation for an MVC controller

Documentation
-------------

http://godoc.org/github.com/ScruffyProdigy/Middleware/controller

Usage
-----

* Create a Controller Class
	* Define a type, it will be a struct
	* You should name the type whatever you want the route to appear as in the URL
		* ie, if your route is "example.com/posts", you should name your type "Posts"
		* if you can't or don't name your type this, define a "RouteName()" function for your type that returns just a string
	* Include in the struct, a controller.Heart
	
*Example:*

		type Posts struct {
			controller.Heart
		}
		
* Define the following method for your type `Find(index string, vars map[string]interface{}) (result Model, found bool)`
	* `index` is the ID of the resource we are looking for
	* `vars` are a list of middleware variables that have been assembled so far
	* `found` is whether or not the resource with the specified ID was found
	* `result` is the resource, if it was found
		* Should have an ID() function that returns index
	* This function will be called to see if a subroute is a valid resource
		* if the request is for "example.com/posts/17", Find("17",vars) will be called
			* If a post "17" exists, you should `return post, true`
			* If not, you should `return nil, false`
			
*Example:*

		type PostModel struct {
			id int
			Author,Title,Content string
		}
		
		func (this *PostModel) ID() string {
			return strconv.Itoa(this.id)
		}

		var posts []*PostModel

		func (this Posts) Find(index string, vars map[string]interface{}) (result controller.Model, found bool) {
			i,err := strconv.Atoi(index)
			if err != nil {
				return nil,false
			}
			i < 0 || i > len(posts) {
				return nil,false
			}
			return posts[i],true
		}
		
* Define operations that are available for your resource by declaring any of the following methods:
	* REST Methods
		* Index() - which is a shortcut for GetCollection()
		* New() - which is a shortcut for GetCollectionNew()
		* Create() - which is a shortcut for PostCollection()
		* Show() - which is a shortcut for GetMember()
		* Edit() - which is a shortcut for GetCollectionEdit()
		* Update() - which is a shortcut for PutMember()
		* Destroy() - which is a shortcut for DeleteMember()
	* Collection/Member Methods
		* Have the form [RequestType], (Collection | Member), [ActionName], "()"
			* If "Collection" is used, it will refer to an action on the collection
			* If "Member" is used, it will refer to an action on a specific member
			* RequestType is optional, and may be "Get", "Post", "Put", or "Delete"
				* The route will only match requests of the chosen type
				* If the RequestType is missing, it will match requests of any type
			* ActionName is optional, and may be any string
				* CollectionStats(), will for instance match /posts/stats
				* MemberStats(), will for instance match /posts/17/stats
				* If ActionName is missing, it will match requests of the type /posts or /posts/17
				
*Example:*

		func (this Posts) Index() {
			this.Set("Posts",posts)
		}
		
		func (this Posts) Create() {
			post := new(PostModel)
			post.Author = this.FormValue("Post[Author]")
			post.Title = this.FormValue("Post[Title]")
			post.Content = this.FormValue("Post[Content]")
			post.id = len(posts)
			posts = append(posts,post)
			this.RespondWith(post)
		}
		
		func (this Posts) Show() {
		}
		
		func (this Posts) Update() {
			post := this.GetVal("Post").(PostModel)
			if author := this.FormValue("Post[Author]");author != "" {
				post.Author = author
			}
			if title := this.FormValue("Post[Title]");title != "" {
				post.Title = title
			}
			if content := this.FormValue("Post[Content]");content != "" {
				post.Content = content
			}
			posts[post.id] = post
		}
		
		func (this Posts) Delete() {
			post := this.GetVal("Post").(PostModel)
			id := post.id
			posts[id] = posts[len(posts)-1]
			posts[id].id = id
			posts = posts[:len(posts)-1]
		}
		
		func (this Posts) GetCollectionFind() {
			found := make([]*PostModel,0,len(posts))
			author := this.FormValue("author")
			title := this.FormValue("title")
			
			for _,post := range(posts) {
				if author != "" && author != post.Author {
					continue
				}
				if title != "" && title != post.Title {
					continue
				}
				found = append(found,post)
			}
			
			this.Set("Posts",found)
		}
		
* Create templates for each GET action that is called
	* Each controller should get it's own folder of templates
		* Each folder should be named the same as the route that is used to invoke the controller
	* See http://golang.com/pkg/text/template for more information on templates
	* All variables set using this.Set() within the control functions will be usable within the templates
	
*Example:*

	/templates/posts/show.tmpl
		
		<h1>{{.Post.Title}}</h1>
		<p>{{.Post.Content}}</p>
		<p><i>by {{.Post.Author}}</i></p>
		
	/templates/posts/index.tmpl
	
		<h1>Posts</h1>
		<ul>
		{{range .Posts}}<li>{{.Title}}</li>
		{{end}}</ul>
		
	/templates/posts/find.tmpl
	
		<h1>Found Posts</h1>
		<ul>
		{{range .Posts}}<li>{{.Title}}</li>
		{{end}}</ul>

* Set up the Rack
	* You will need to set up the routes at the same time
		* To get the route for the Resource, call NewResource() with an instance of your controller class
	* The Templater Middleware is required to be used before this Middleware
		* Each controller should have a separate folder of templates; you must send Templater the folder that contains each of those folders
	* Other Middleware may be recommended to be used in tandem

*Example:*
	
		func main() {
			root := controller.NewRoot()
			postsRoute := controller.NewResource(Posts)
			root.AddRoute(postsRoute)
			
			rackup := rack.New()
			rackup.Add(templater.GetTemplates("./templates"))
			rackup.Add(root)
			
			conn := httper.HttpConnection(":5001")
			go conn.Go(rackup)
		}