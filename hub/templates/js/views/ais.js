window.AIsView = Backbone.View.extend({
	
	template: _.template($('#ais_underscore').html()),

	events: {
	  'click .add-ai button': 'createNewAI',
		'click .delete-button': 'deleteAI',
	},

	initialize: function() {
	  this.collection = new AIs();
		this.listenTo(this.collection, 'reset', this.render);
		this.listenTo(this.collection, 'add', this.render);
		this.listenTo(this.collection, 'remove', this.render);
		this.collection.fetch({ reset: true });
	},

	createNewAI: function(ev) {
	  ev.preventDefault();
		this.collection.create({
		  Name: $('.new-ai-name').val(),
			URL: $('.new-ai-url').val(),
			IsOwner: true,
		});
	},

	deleteAI: function(ev) {
	  ev.preventDefault();
	  var ai = this.collection.get($(ev.target).attr('data-id'));
		ai.destroy();
	},

  render: function() {
		var that = this;
    that.$el.html(that.template({}));
		that.collection.each(function(ai) {
		  if (ai.get('IsOwner')) {
				that.$('table').append('<tr><td>' + ai.get('Name') + '</td><td>' + ai.get('URL') + '</td><td><button data-id="' + ai.get('Id') + '" class="btn btn-xs delete-button">Delete</button></a></td></tr>');
			} else {
				that.$('table').append('<tr><td>' + ai.get('Name') + '</td><td>' + ai.get('URL') + '</td><td>' + ai.get('Owner') + '</td></tr>');
			}
		});
		if (window.session.user.loggedIn()) {
		  that.$('.add-ai').show();
		} else {
		  that.$('.add-ai').hide();
		}
		return that;
	},

});
