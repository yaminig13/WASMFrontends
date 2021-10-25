module Prism
  @@instances = {}

  def self.instances
    @@instances
  end

  def self.mount(component)
    Mount.new(component)
  end

  class Mount

    def initialize(component)
      @component = component
    end

    def render
      JSON::stringify(@component.render)
    end

    def dispatch(messageJSON)
      message = JSON::parse(messageJSON)

      instance = Prism.instances[message["instance"]]

      instance.process(message)
    end

    def event(eventJSON, id)
      DOM.event(JSON::parse(eventJSON), id)
    end

    def display_location(eventJSON, id)
      Location.event(JSON::parse(eventJSON), id)
    end

    def http_response(responseJSON, id)
      HTTP._response(HTTP::Response.from_json(responseJSON), id)
    end

    def cookie_value(cookieValues,id)
      Storage.storage(cookieValues,id)
    end

    def session_value(sessionValue,id)
      Storage.storage(sessionValue,id)
    end

    def local_value(localValue,id)
      Storage.storage(localValue,id)
    end

    def send_msg(serverMsg,id)
      WebSocket.displayMsg(serverMsg,id)
    end

    def display_value(valueData,id)
      DOM.event(valueData,id)
    end

    def param_value(valueData,id)
      Navigation.output(valueData,id)
    end

    def indexeddb_response(valueData,id)
      IndexedDB.displayMsg(valueData,id)
    end
    def worker_fn(eventData,id)
      WorkerThread._response(eventData,id)
    end
    def camera_output(eventData,id)
      CameraAPI.event(eventData,id)
    end
   

  end

  class Component
    TAG_NAMES = %i[
      a
      abbr
      address
      area
      article
      aside
      audio
      b
      base
      bdi
      bdo
      blockquote
      body
      br
      button
      canvas
      caption
      cite
      code
      col
      colgroup
      dd
      del
      details
      dfn
      dir
      div
      dl
      dt
      em
      embed
      fieldset
      figcaption
      figure
      footer
      form
      h1
      h2
      h3
      h4
      h5
      h6
      head
      header
      hgroup
      hr
      html
      i
      iframe
      img
      input
      ins
      kbd
      keygen
      label
      legend
      li
      link
      main
      map
      mark
      menu
      meta
      nav
      noscript
      object
      ol
      optgroup
      option
      p
      param
      pre
      progress
      q
      rp
      rt
      ruby
      s
      samp
      script
      section
      select
      small
      source
      span
      strong
      style
      sub
      summary
      sup
      table
      tbody
      td
      textarea
      tfoot
      th
      thead
      time
      title
      tr
      u
      ul
      video
    ]

    TAG_NAMES.each do |tag|
      define_method(tag) do |*args|
        options = {}
        className = ""
        children = []

        result = {}

        until args.empty?
          arg = args.shift

          case arg
          when String
            if arg.start_with?(".")
              className = arg
            else
              children = [text(arg)]
            end
          when Array
            children = arg
          when Object
            options = arg
          end
        end

        options.each do |key, value|
          next if value.is_a?(EventHandler) || (key.to_s).chars.first == '_'
          result[key] = value
        end

        result[:_type] = tag
        result[:_class] = className
        options[:_children] = children || []

        result[:_children] = options[:_children].compact.map do |child|
          if child.is_a?(Prism::Component)
            child.render
          elsif child.is_a?(String)
            text(child)
          else
            child
          end
        end

        result[:on] ||= {}

        options.each do |key, value|
          key_as_string = key.to_s
          next unless key_as_string.start_with?('on')
          event_name = key_as_string.sub('on', '').downcase
          result[:on][event_name] = _register_handler(value)
        end

        if options[:on]
          event_handlers = {}

          options[:on].each do |key, value|
            event_handlers[key] = value.to_hash
          end

          result[:on] = event_handlers
        end

        result
      end
    end

    def text(t)
      {:type => "text", :content => t.to_s}
    end

    def call(method_name, *args)
      EventHandler.new(self, method_name).with(*args)
    end

    def stop_propagation
      EventHandler.new(self, method_name).stop_propagation
    end

    def prevent_default
      EventHandler.new(self, method_name).prevent_default
    end

    def _register_handler(handler)
      Prism.instances[handler.id] = handler # TODO - this is a memory leak
      handler.to_hash
    end

    def render
      raise "Unimplemented render method for #{self.class.name}"
    end
  end

  class EventHandler
    attr_reader :method_name, :id

    def initialize(instance, method_name, args = [], options = {})
      @instance = instance
      @id = self.object_id
      @method_name = method_name
      @args = args
      @options = {prevent_default: false, stop_propagation: false}.merge(options)
    end

    def with(*additionalArgs)
      new_args = additionalArgs.map { |a| {type: :constant, value: a} }

      EventHandler.new(@instance, method_name, @args + new_args, @options)
    end

    def with_event
      EventHandler.new(@instance, method_name, @args + [{type: :event}], @options)
    end

    def with_event_data(*property_names)
      new_args = property_names.map { |item| {type: :event_data, key: item } }

      EventHandler.new(@instance, method_name, @args + new_args, @options)
    end

    def with_target_data(*items)
      target_args = items.map { |item| {type: :target_data, key: item } }
      EventHandler.new(@instance, method_name, @args + target_args, @options)
    end

    def prevent_default
      EventHandler.new(@instance, method_name, @args, @options.merge(prevent_default: true))
    end

    def stop_propagation
      EventHandler.new(@instance, method_name, @args, @options.merge(stop_propagation: true))
    end

    def process(message)
      call_args = []
      message_args = message["args"]

      @args.each do |arg|
        case arg[:type]
        when :constant
          call_args << arg[:value]
        else
          call_args << message_args.shift
        end
      end

      if @instance.respond_to? message["type"]
        @instance.send(message["type"], *call_args)
      else
        raise "Component #{@instance.class} has no method ##{message["type"]}"
      end
    end

    def to_hash
      {
        instance: @id,
        type: @method_name,
        args: @args,
        prevent_default: @options[:prevent_default],
        stop_propagation: @options[:stop_propagation]
      }
    end
  end
end

module DOM
  @@event_id = 0
  @@listeners = {}

  def self.get_event_id
    @@event_id += 1

    @@event_id.to_s
  end

  def self.add_listener(id, &block)
    @@listeners[id] = block
  end

  def self.select(selector)
    ElementSelection.new(selector)
  end

  def self.event(eventData, id)
    @@listeners[id].call(eventData)
  end

  class ElementSelection
    def initialize(selector)
      @selector = selector
    end

    # def getElement(&block)
    #   id = DOM.get_event_id
    #   DOM.add_listener(id, &block)
    #   InternalDOM.get_element(@selector, id)
      
    # end

    def on(eventName, &block)
      id = DOM.get_event_id
      InternalDOM.add_event_listener(@selector, eventName, id)
      DOM.add_listener(id, &block)
    end

    def getValue(&block)
      id = DOM.get_event_id
      DOM.add_listener(id, &block)
      InternalDOM.get_value(@selector, id)
      
    end
  end
end

module CameraAPI
  @@event_id = 0
  @@listeners = {}

  def self.get_event_id
    @@event_id += 1

    @@event_id.to_s
  end

  def self.add_listener(id, &block)
    @@listeners[id] = block
  end

  def self.event(eventData, id)
    @@listeners[id].call(eventData)
  end

  def self.enable_camera(&block)
    id = CameraAPI.get_event_id
    CameraAPI.add_listener(id, &block)
    InternalCameraAPI.enable_camera(id)
  end
end

module Location
  @@event_id = 0
  @@listeners = {}

  def self.get_event_id
    @@event_id += 1

    @@event_id.to_s
  end

  def self.add_listener(id, &block)
    @@listeners[id] = block
  end

  def self.event(eventData, id)
    @@listeners[id].call(eventData)
  end

  def self.find(&block)
      id = Location.get_event_id
      InternalLocation.get_location(id)
      Location.add_listener(id, &block)
    end

end

module WebSocket
  
  @@event_id = 0
  @@listeners = {}

  def self.get_event_id
    @@event_id += 1

    @@event_id.to_s
  end

  def self.add_listener(id, &block)
    @@listeners[id] = block
  end

  def self.displayMsg(serverMsg,id)
    @@listeners[id].call(serverMsg)

  end
  def self.openConnection(url,clientMsg,&block)
    id=WebSocket.get_event_id
    InternalWebSocket.open_connection(url,clientMsg,id)
    WebSocket.add_listener(id,&block)
  end
end

module PushNotifAPI
  def self.sendPushNotif()
    InternalPushNotifAPI.send_pushNotif()
  end
end

module IndexedDB

  @@event_id = 0
  @@listeners = {}

  def self.get_event_id
    @@event_id += 1

    @@event_id.to_s
  end

  def self.add_listener(id, &block)
    @@listeners[id] = block
  end

  def self.displayMsg(response,id)
    @@listeners[id].call(response)

  end

  def self.createDB(name,items)
    InternalIndexedDB.create_db(name,items)
  end

  def self.openDB(name,term,&block)
    id=IndexedDB.get_event_id
    InternalIndexedDB.open_db(name,term,id)
    IndexedDB.add_listener(id,&block)
  end
end

module Storage

  @@event_id = 0
  @@listeners = {}

  def self.get_event_id
    @@event_id += 1

    @@event_id.to_s
  end

  def self.add_listener(id, &block)

    @@listeners[id] = block

  end

  def self.storage(storageData, id)
    @@listeners[id].call(storageData)
  end

  def self.type(storageType)
    StorageType.new(storageType)
  end

  class StorageType
    def initialize(storageType)
      @storageType=storageType
    end

    def set(name,value)
      if @storageType=="LocalStorage"
        InternalLocal.set_local(name,value)
      elsif @storageType=="SessionStorage"
        InternalSession.set_session(name,value)
      elsif @storageType=="Cookie"
        InternalCookie.set_cookie(value)
      else
        puts "error"
      end
    end

    def get(name, &block)
      
      id = Storage.get_event_id
      Storage.add_listener(id, &block)

      if @storageType=="LocalStorage"
        InternalLocal.get_local(name,id)
      elsif @storageType=="SessionStorage"
        InternalSession.get_session(name,id)
      elsif @storageType=="Cookie"
        InternalCookie.get_cookie(id)
      else
        puts "error"
      end

    end
  end
end

module Navigation

  @@event_id = 0
  @@listeners = {}
  

  def self.get_event_id
    @@event_id += 1

    @@event_id.to_s
  end

  def self.add_listener(id, &block)
    @@listeners[id] = block
  end

  def self.navigateto(url,customId)
    InternalNavigation.navigate_to(url,customId)
  end

  def self.output(outputData,id)
     @@listeners[id].call(outputData)
  end

  def self.getParam(pname,&block)
    id=Navigation.get_event_id
    Navigation.add_listener(id,&block)

    InternalNavigation.get_param(pname,id)
  end

end

module VibrationAPI
  def self.vibrate(time)
    InternalVibrationAPI.vibrate_for(time)
  end
end

module NotificationAPI
  def self.sendNotif(title,body)
    InternalNotificationAPI.send_notif(title,body)
  end
end

module Alert
  def self.showAlert(msg)
    InternalAlert.show_alert(msg)
  end
end

module WorkerThread
   @@event_id = 0
  @@listeners = {}

  def self.get_event_id
    @@event_id += 1

    @@event_id.to_s
  end

  def self.add_listener(id, &block)
    @@listeners[id] = block
  end

  def self._response(text,id)
    @@listeners[id].call(text)
  end

  def self.receiveMsg(worker,&block)
    id = WorkerThread.get_event_id
    InternalWorkerThread.receive_msg(worker,id)
    WorkerThread.add_listener(id, &block)

  end

   def self.startWorker(url,&block)
    id = WorkerThread.get_event_id
    WorkerThread.add_listener(id, &block)

    InternalWorkerThread.start_worker(url,id)

  end
end

module HTTP
  @@event_id = 0
  @@listeners = {}

  def self.get_event_id
    @@event_id += 1

    @@event_id.to_s
  end

  def self.add_listener(id, &block)
    @@listeners[id] = block
  end

  def self.get(url, &block)
    id = HTTP.get_event_id

    InternalHTTP.http_request(url, id)

    HTTP.add_listener(id, &block)
  end

  def self.post(url,data)
    InternalHTTP.http_post(url, data)
  end

  def self._response(text, id)
    @@listeners[id].call(text)
  end

  class Response < Struct.new(:body)
    def self.from_json(json)
      data = JSON::parse(json)

      new(data["body"])
    end
  end
end

CONFIG = JSON::parse(JSON_CONFIG)
