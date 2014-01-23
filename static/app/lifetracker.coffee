$(document).ready ->
  app = app || window

  app.Router = Backbone.Router.extend
    routes:
      "": "loadBoard"
      "calendar": "loadCalendar"
      "board": "loadBoard"

    hideAll: ->
      app.calendarView.hide()
      app.accomplishmentsView.hide()

    loadCalendar: ->
      @hideAll()
      app.calendarView.render()

    loadBoard: ->
      @hideAll()
      app.accomplishmentsView.render()

  app.router = new app.Router
  Backbone.history.start();