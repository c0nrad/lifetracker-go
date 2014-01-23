$(document).ready ->
  app = window || {}


  # Returns relative to cellblock number
  # If jan 1 is a Tuesday then 
  app.CalendarModel = Backbone.Model.extend
    defaults: ->
      currMonth: moment().format('MMMM')
      currYear: moment().format('YYYY')
      currDay: moment().format('DD')

    firstCellBlock: ->
      days = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"]
      firstDay = "#{@get('currMonth')} 1 #{@get('currYear')}"
      dayName = moment(firstDay).format('dddd')
      days.indexOf(dayName)


  app.CalendarCellView = Backbone.View.extend
    template: _.template($('#calendarCell-template').html())

    initialize: ->

    render: ->
      $(@el).html @template
        day: @model.day
        body: @model.body
      @

  app.CalendarView = Backbone.View.extend
    el: $('#calendarView')

    initialize: ->
      @model = new app.CalendarModel()

    hide: ->
      $(@el).hide()

    unhide: ->
      $(@el).show()

    render: -> 
      $(@el).find("#calendarTitle").text("#{@model.get('currMonth')} #{@model.get('currYear')}")
      firstDay = moment("#{@model.get('currMonth')} 1 #{@model.get('currYear')}")

      currDay = firstDay
      firstCellNumber = @model.firstCellBlock()

      for x in [0..firstDay.daysInMonth()]
        cellView = new app.CalendarCellView {model: {day: x}}

        currAccomplishment = app.accomplishments.getAccomplishment(currDay)
        
        if currAccomplishment?
          cellView = new app.CalendarCellView {model: {day: x, body: currAccomplishment.get('Body')}}
        else
          cellView = new app.CalendarCellView {model: {day: x, body: ""}}

        $(@el).find("#cell"+(x+firstCellNumber) ).html( cellView.render().el )
        currDay.add('days',1)
      @unhide()

  app.calendarView = new app.CalendarView