# Promotion system

<img src="docs/images/Cake-mobile.svg" alt="all component"/>

## Problem description
**Build a promotion system to support 100 first login user per campaign when client register a new account on Cake system. The 100 first login users will get 30% discount voucher when they top-up the mobile phone‘s fee (Money transfer from bank account) via Cake app.**
## Analyze
### Functional requirement
- Manipulate campaign, rule and benefit for users.
    - Rule:
        - User is new account on system.
        - Limit 100 first user only.
        - User do top-up mobile phone action.
    - Effect:
        - Discount 30%
- The effect will be execute only after user has finished payment action.
- End user will action via Cake mobile apps.
### Non-functional requirement
- System can support >= 100,000 concurrent users on the campaign:
    - The system should be scalable.
## Component definition

- **Rule**: A series of condition define by promotion teams. If the conditions match, the **effect** will be executed. In simple language, this is **if statement**. Eg: "If this user is in 100 first log in users and user top-up the mobile phone".
- **Effect**: The **effect** after user finish an action that met  **rule** dragged on it. This is the **then** statement in **if-then** clause. Eg: "discount 30%".
- **Rule-engine**: Solve the rule-effect relationship by do **rule** validation, then do the **effect**.
- **Offers**: Store the rule-effect metadata and pass needed input to the **rule-engine**.
- **Campaigns**: Define by promotion team, the campaign contain offers. In other word, the offers will be appeared when the campaign show it.
- **Invalid cache**: Store the invalidate user-campaign-offer metadata.
- **user-campaign-offer metadata**: The composite key define the association from user data, campaign and offer metadata. This can be use for some meanings: The user can or cannot use the offer from the campaign.
- **Validate campaign**: Do validation of  campaign-offer key. This can be use to update is campaign is still activating in realtime.
- **Validation**: Do validation on user-campaign-offer.
- **User Event**: An pipeline queue that received and keep event from user, including payment action.
- **Payment**: Do necessary step when user finisher paying action. This's subscribe to UserEvent queue to get the input data.
- **User:** Service to handle user information.
- **Pro-processed database**: The data will be pre-calculator first and try to push them to this database. This will improve time for validation and suggest offers for users. This should be a READ optimization database
- **Data processing**:  Handle the data pre-calculation and push the centralize database.
- **User's balance**: Do the money balancing for user.
- **Cake Gate**: Internet public gateway received and routing user event from Cake mobile apps directly. This include authen/author/rate limiter also.
- **API Gateway** or **Payment Gateway**: Payment internal gateway route to needed services for payment and promotion business logic.

## Campaign creation side

<img src="docs/images/Cake-campaign-creator.svg">

When promotion team create/update campaign, they will create or select existed rule and effect also. Some will invalid or out-of-date, so for fast validation, we push them into "invalid-cache" first.
Right after campaign is created or new non-payment event from user is received, the system will calculate to answer the question: "For each selected users, which campaigns and offers is available to them". The data will be pre-calculate in this phase and will be store in centralized database that optimize for READ action. This will help to improve response time when user try to get list of theirs offers and do payment action. The data will be store with the composite key from `userid-campaignid-offersid` .

## When user list their offers

<img src="docs/images/Cake-get-offers.svg">

When user try to get list of offers, the simple flow can be:
- System try to find all valid offers from **preprocessed database** first.
- The offers is in invalid cache data will be removed.

##  When user finisher their payment action.

<img src="docs/images/Cake-issue-offer.svg">

- Each user's payment event will be pushed into **UserEvent** queue.
- The subscribe **payment** service will take the event with input from user, try to get **userdata** from **User** service.
- The validation step will be execute as when user try to get list of offers again.
- When the **validation** step is ok, the offer's effect will be execute. User's balance will be calculate and return success response when the all bank's transaction is succeed.

### User's balance calculation

<img src="docs/images/Cake-balance.svg">

Example: When user top-up 100k VND to mobile apps:
- System will notify user need 70K VND in balance is enough if user use the new member voucher.
- When use finish payment action from the apps:
    - System will execute money transferring from **Campaign Bank account** to **Merchant** first.
    - System will execute money transferring from **User Bank Acocunt** to **Campaign Bank account**.
    - Update calculated balance into database and return response to user.
