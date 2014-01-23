$(document).ready ->
  app = window || {}

  app.Accomplishment = Backbone.Model.extend
    defaults: ->
      Name: "Anonymous"
      Body: "I worked on calendar view...."
      Date: new Date()
      ImagePath: ""

  app.Accomplishments = Backbone.Collection.extend
    model: app.Accomplishment

    getAccomplishment: (date) ->
      @find (model) ->
        a = moment(model.get('Date'))
        b = moment(date)
        a.month() == b.month() and b.year() == a.year() and a.dayOfYear() == b.dayOfYear() 

  app.AccomplishmentView = Backbone.View.extend
    tagName: "div"
    className: "accomplishment"
    template: _.template($('#accomplishment-template').html())

    initialize: ->
      # _.bindAll @, 'render'

    render: ->
      $(@el).html @template (@model.toJSON())
      $(@el).find('.date').text moment(@model.get('Date')).format("MMM D")
      @

  app.AccomplishmentsView = Backbone.View.extend
    el: $("#publicView")

    initialize: ->
      @collection = app.accomplishments = new app.Accomplishments app.accomplishmentData

    hide: ->
      $(@el).hide()

    unhide: ->
      $(@el).show()

    render: ->
      $(@el).find("#accomplishments").html("")
      @collection.each (model) =>
        accomplishmentView = new app.AccomplishmentView {model: model}
        $(@el).find("#accomplishments").append accomplishmentView.render().el 

      container = document.querySelector('#accomplishments');
      imagesLoaded container, ->
        msnry = new Masonry container, 
          columnWidth: 75,
          itemSelector: '.item'
      @unhide()


  app.accomplishmentsView = new app.AccomplishmentsView

  