package main

import (
	"html/template"
	//	"strconv"
	"time"
)

var (
	hours    []int
	minutes  []int
	days     = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	actions  = []string{"fadeup", "fadedown"}
	MainForm = template.New("MainForm")
)

func init() {
	hours = make([]int, 24)
	for h := 0; h < 24; h++ {
		// hours[h] = strconv.FormatInt(int64(h), 10)
		hours[h] = h
	}
	minutes = make([]int, 60)
	for m := 0; m < 60; m++ {
		minutes[m] = m
	}
}

type FormData struct {
	CurrentTime time.Time
	Hours       []int
	Minutes     []int
	Actions     []string
	Timers      map[string]Timer
}

const webform = `
<html>
<head>
</head>
<body>
  <fieldset><legend>Current Time</legend>
  <form name="manual" method="POST">
    <table>
    <tr><td>{{.CurrentTime}}</tr></td>
    <tr><td>
    <label for="level">Level</label>
      <input type="number" id="level" name="level" min="0" max="100">
      <input name="SubmitAction" type="submit" value="Set">
    </td></tr>
    </table>
  </form>
  </fieldset>
  <fieldset><legend>Time Programs</legend>

  {{ $hourList := .Hours }}
  {{ $minuteList := .Minutes }}
  {{ $actionList := .Actions }}
  
  {{range $key, $value := .Timers}}
    <fieldset><legend>{{$key}}</legend>
    <form name="{{$key}}" method="POST">
    <input type="hidden" name="timer" value="{{$key}}"/>
    <table>
    <tr><td>
    <label for="enabled">Enabled</label>
    {{if $value.Enabled}}
      <input type="checkbox" name="enabled" id="enabled" checked>
    {{else}}
      <input type="checkbox" name="enabled" id="enabled">
    {{end}}

    <label for="hour">Hour</label>
    <select name="hour" id="hour">
      {{ range $hourList}}
        {{if eq $value.Hour . }}
          <option value="{{.}}" selected>{{.}}</option>
        {{else}}
          <option value="{{.}}">{{.}}</option>
        {{end}}
      {{ end }}
    </select>

    <label for="minute">Minute</label>
    <select name="minute" id="minute">
      {{ range $minuteList}}
        {{if eq $value.Minute . }}
          <option value="{{.}}" selected>{{.}}</option>
        {{else}}
          <option value="{{.}}">{{.}}</option>
        {{end}}
      {{ end }}
    </select>
    
    </td>
    </tr>
    <tr>
    <td>
    <label for="Sunday">Sunday</label>
    {{if $value.Sunday }}
      <input type="checkbox" name="Sunday" id="Sunday" checked>
    {{else}}
      <input type="checkbox" name="Sunday" id="Sunday">
    {{end}}

    <label for="Monday">Monday</label>
    {{if $value.Monday }}
      <input type="checkbox" name="Monday" id="Monday" checked>
    {{else}}
      <input type="checkbox" name="Monday" id="Monday">
    {{end}}
      
    <label for="Tuesday">Tuesday</label>
    {{if $value.Tuesday }}
      <input type="checkbox" name="Tuesday" id="Tuesday" checked>
    {{else}}
      <input type="checkbox" name="Tuesday" id="Tuesday">
    {{end}}
      
    <label for="Wednesday">Wednesday</label>
    {{if $value.Wednesday }}
      <input type="checkbox" name="Wednesday" id="Wednesday" checked>
    {{else}}
      <input type="checkbox" name="Wednesday" id="Wednesday">
    {{end}}
      
    <label for="Thursday">Thursday</label>
    {{if $value.Thursday }}
      <input type="checkbox" name="Thursday" id="Thursday" checked>
    {{else}}
      <input type="checkbox" name="Thursday" id="Thursday">
    {{end}}
      
    <label for="Friday">Friday</label>
    {{if $value.Friday }}
      <input type="checkbox" name="Friday" id="Friday" checked>
    {{else}}
      <input type="checkbox" name="Friday" id="Friday">
    {{end}}
      
    <label for="Saturday">Saturday</label>
    {{if $value.Saturday }}
      <input type="checkbox" name="Saturday" id="Saturday" checked>
    {{else}}
      <input type="checkbox" name="Saturday" id="Saturday">
    {{end}}
    </td>
    </tr>
    <tr>
    <td>
    <label for="Action">Action</label>
    <select name="Action" id="Action">
    {{ range $actionList}}
      {{if eq $value.Action . }}
        <option value="{{.}}" selected>{{.}}</option>
      {{else}}
        <option value="{{.}}">{{.}}</option>
      {{end}}
    {{ end }}
    </select>
    </td>
    </tr>
    <tr>
    <td>
    <input name="SubmitAction" type="submit" value="Update">
    <input name="SubmitAction" type="submit" value="Delete">
    </td>
    </tr>
    </table>
    </form>
    </fieldset>
  {{ end }}
  
  <fieldset><legend>New Timer</legend>
  <form name="NewTimer" method="POST">
  <input type="text" name="timer" value=""/>
  <table>
  <tr><td>
  <label for="enabled">Enabled</label>
    <input type="checkbox" name="enabled" id="enabled">

  <label for="hour">Hour</label>
  <select name="hour" id="hour">
    {{ range $hourList}}
        <option value="{{.}}">{{.}}</option>
    {{ end }}
  </select>

  <label for="minute">Minute</label>
  <select name="minute" id="minute">
    {{ range $minuteList}}
        <option value="{{.}}">{{.}}</option>
    {{ end }}
  </select>
  
  </td>
  </tr>
  <tr>
  <td>
  <label for="Sunday">Sunday</label>
  <input type="checkbox" name="Sunday" id="Sunday">


  <label for="Monday">Monday</label>
  <input type="checkbox" name="Monday" id="Monday">
    
  <label for="Tuesday">Tuesday</label>
  <input type="checkbox" name="Tuesday" id="Tuesday">
    
  <label for="Wednesday">Wednesday</label>
  <input type="checkbox" name="Wednesday" id="Wednesday">

  <label for="Thursday">Thursday</label>
  <input type="checkbox" name="Thursday" id="Thursday">
    
  <label for="Friday">Friday</label>
  <input type="checkbox" name="Friday" id="Friday">
    
  <label for="Saturday">Saturday</label>
  <input type="checkbox" name="Saturday" id="Saturday">
  </td>
  </tr>
  <tr>
  <td>
  <label for="Action">Action</label>
  <select name="Action" id="Action">
  {{ range $actionList}}
    <option value="{{.}}">{{.}}</option>
  {{ end }}
  </select>
  </td>
  </tr>
  <tr>
  <td>
  <input name="SubmitAction" type="submit" value="Create">
  </td>
  </tr>
  </table>
  </form>
  </fieldset>
  
  </fieldset>
</body>
</html>
`
