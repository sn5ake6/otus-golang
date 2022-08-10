# file: features/calendar.feature

# http://calendar:8888/

Feature: Calendar
	In order to test application behavior
	As API client of calendar service
	I need to be able to run features

	Scenario: Service is available
		When I send "GET" request to "http://calendar:8888/events/"
		Then The response code should be 200

	Scenario: Event is created
		When I send "POST" request to "http://calendar:8888/events/" with "application/json" data:
		"""
		{
			"id":"512bc5cd-01e9-4639-99a2-d42fe25dec62",
			"title":"some title",
			"beginAt":"2022-07-24T16:00:00Z",
			"endAt":"2022-07-24T18:00:00Z",
			"description":"some description",
			"userId":"6b216e09-7ab3-41f9-ba57-cc94d45fe759",
			"notifyAt":"2022-07-24T15:00:00Z"
		}
		"""
		Then The response code should be 201

	Scenario: Duplicate event is not created
		When I send "POST" request to "http://calendar:8888/events/" with "application/json" data:
		"""
		{
			"id":"512bc5cd-01e9-4639-99a2-d42fe25dec62",
			"title":"some title",
			"beginAt":"2022-07-24T16:00:00Z",
			"endAt":"2022-07-24T18:00:00Z",
			"description":"some description",
			"userId":"6b216e09-7ab3-41f9-ba57-cc94d45fe759",
			"notifyAt":"2022-07-24T15:00:00Z"
		}
		"""
		Then The response code should not be 201

	Scenario: Created event is received
		When I send "GET" request to "http://calendar:8888/events/512bc5cd-01e9-4639-99a2-d42fe25dec62/"
		Then The response code should be 200
		And The response should match text:
		"""
		{
			"id":"512bc5cd-01e9-4639-99a2-d42fe25dec62",
			"title":"some title",
			"beginAt":"2022-07-24T16:00:00Z",
			"endAt":"2022-07-24T18:00:00Z",
			"description":"some description",
			"userId":"6b216e09-7ab3-41f9-ba57-cc94d45fe759",
			"notifyAt":"2022-07-24T15:00:00Z",
			"notifiedAt":"0001-01-01T00:00:00Z"
		}
		"""

	Scenario: Event is updated
		When I send "PUT" request to "http://calendar:8888/events/512bc5cd-01e9-4639-99a2-d42fe25dec62/" with "application/json" data:
		"""
		{
			"id":"512bc5cd-01e9-4639-99a2-d42fe25dec62",
			"title":"updated title",
			"beginAt":"2022-07-25T16:00:00Z",
			"endAt":"2022-07-25T18:00:00Z",
			"description":"updated description",
			"userId":"6b216e09-7ab3-41f9-ba57-cc94d45fe759",
			"notifyAt":"2022-07-25T15:00:00Z"
		}
		"""
		Then The response code should be 200

	Scenario: Updated event is received
		When I send "GET" request to "http://calendar:8888/events/512bc5cd-01e9-4639-99a2-d42fe25dec62/"
		Then The response code should be 200
		And The response should match text:
		"""
		{
			"id":"512bc5cd-01e9-4639-99a2-d42fe25dec62",
			"title":"updated title",
			"beginAt":"2022-07-25T16:00:00Z",
			"endAt":"2022-07-25T18:00:00Z",
			"description":"updated description",
			"userId":"6b216e09-7ab3-41f9-ba57-cc94d45fe759",
			"notifyAt":"2022-07-25T15:00:00Z",
			"notifiedAt":"0001-01-01T00:00:00Z"
		}
		"""

	Scenario: Not empty events on day are received
		When I send "GET" request to "http://calendar:8888/events/?period=day&date=2022-07-25%2010:00:00"
		Then The response code should be 200
		And The response should match text:
		"""
		[
			{
				"id":"512bc5cd-01e9-4639-99a2-d42fe25dec62",
				"title":"updated title",
				"beginAt":"2022-07-25T16:00:00Z",
				"endAt":"2022-07-25T18:00:00Z",
				"description":"updated description",
				"userId":"6b216e09-7ab3-41f9-ba57-cc94d45fe759",
				"notifyAt":"2022-07-25T15:00:00Z",
				"notifiedAt":"0001-01-01T00:00:00Z"
			}
		]
		"""

	Scenario: Empty events on day are received
		When I send "GET" request to "http://calendar:8888/events/?period=day&date=2022-07-20%2010:00:00"
		Then The response code should be 200
		And The response should match text:
		"""
		[]
		"""

	Scenario: Not empty events on week are received
		When I send "GET" request to "http://calendar:8888/events/?period=week&date=2022-07-25%2010:00:00"
		Then The response code should be 200
		And The response should match text:
		"""
		[
			{
				"id":"512bc5cd-01e9-4639-99a2-d42fe25dec62",
				"title":"updated title",
				"beginAt":"2022-07-25T16:00:00Z",
				"endAt":"2022-07-25T18:00:00Z",
				"description":"updated description",
				"userId":"6b216e09-7ab3-41f9-ba57-cc94d45fe759",
				"notifyAt":"2022-07-25T15:00:00Z",
				"notifiedAt":"0001-01-01T00:00:00Z"
			}
		]
		"""

	Scenario: Empty events on week are received
		When I send "GET" request to "http://calendar:8888/events/?period=week&date=2022-07-10%2010:00:00"
		Then The response code should be 200
		And The response should match text:
		"""
		[]
		"""

	Scenario: Not empty events on month are received
		When I send "GET" request to "http://calendar:8888/events/?period=month&date=2022-07-25%2010:00:00"
		Then The response code should be 200
		And The response should match text:
		"""
		[
			{
				"id":"512bc5cd-01e9-4639-99a2-d42fe25dec62",
				"title":"updated title",
				"beginAt":"2022-07-25T16:00:00Z",
				"endAt":"2022-07-25T18:00:00Z",
				"description":"updated description",
				"userId":"6b216e09-7ab3-41f9-ba57-cc94d45fe759",
				"notifyAt":"2022-07-25T15:00:00Z",
				"notifiedAt":"0001-01-01T00:00:00Z"
			}
		]
		"""

	Scenario: Empty events on month are received
		When I send "GET" request to "http://calendar:8888/events/?period=month&date=2022-06-20%2010:00:00"
		Then The response code should be 200
		And The response should match text:
		"""
		[]
		"""

	Scenario: Old event is created
		When I send "POST" request to "http://calendar:8888/events/" with "application/json" data:
		"""
		{
			"id":"6b216e09-7ab3-41f9-ba57-cc94d45fe759",
			"title":"old event title",
			"beginAt":"2021-07-24T16:00:00Z",
			"endAt":"2021-07-24T18:00:00Z",
			"description":"old event description",
			"userId":"6b216e09-7ab3-41f9-ba57-cc94d45fe759",
			"notifyAt":"2021-07-24T15:00:00Z",
			"notifiedAt":"2021-07-24T15:00:01Z"
		}
		"""
		Then The response code should be 201

	Scenario: Old event is received
		When I send "GET" request to "http://calendar:8888/events/6b216e09-7ab3-41f9-ba57-cc94d45fe759/"
		Then The response code should be 200
		And The response should match text:
		"""
		{
			"id":"6b216e09-7ab3-41f9-ba57-cc94d45fe759",
			"title":"old event title",
			"beginAt":"2021-07-24T16:00:00Z",
			"endAt":"2021-07-24T18:00:00Z",
			"description":"old event description",
			"userId":"6b216e09-7ab3-41f9-ba57-cc94d45fe759",
			"notifyAt":"2021-07-24T15:00:00Z",
			"notifiedAt":"2021-07-24T15:00:01Z"
		}
		"""

	Scenario: Event is notified
		When I wait "1m" and send "GET" request to "http://calendar:8888/events/512bc5cd-01e9-4639-99a2-d42fe25dec62/"
		Then The response code should be 200
		And The response should not contain text:
		"""
		"notifiedAt":"0001-01-01T00:00:00Z"
		"""

	Scenario: Old event is cleared
		When I send "GET" request to "http://calendar:8888/events/6b216e09-7ab3-41f9-ba57-cc94d45fe759/"
		Then The response code should not be 200

	Scenario: Calendar event is deleted
		When I send "DELETE" request to "http://calendar:8888/events/512bc5cd-01e9-4639-99a2-d42fe25dec62/"
		Then The response code should be 200

	Scenario: Deleted calendar event is not received
		When I send "GET" request to "http://calendar:8888/events/512bc5cd-01e9-4639-99a2-d42fe25dec62/"
		Then The response code should not be 200