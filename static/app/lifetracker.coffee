$(document).ready ->
  app = app || window

  app.Router = Backbone.Router.extend
    routes:
      "": "loadPublic"
      "calendar": "loadCalendar"
      "public": "loadPublic"

    hideAll: ->
      app.calendarView.hide()
      app.accomplishmentsView.hide()

    loadCalendar: ->
      @hideAll()
      app.calendarView.render()
      console.log "LOAD CALENDAR MUTHERFUCKER"

    loadPublic: ->
      @hideAll()
      app.accomplishmentsView.render()
      console.log "LOAD PUBLIC MUTHERFUCKER"

  app.router = new app.Router
  Backbone.history.start();