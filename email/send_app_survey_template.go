package sendemail

const AppSurveyTemplate = `
<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<meta http-equiv="X-UA-Compatible" content="IE=edge" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{.Subject}}</title>
	
	<style type="text/css">
	@import url('https://fonts.googleapis.com/css2?family=Afacad:wght@400;500;600;700&display=swap'); /* Added Afacad */
	@import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap');
	
	</style>
</head>
<body style="margin: 0px;">

	<center class="wrapper" style="width: 100%; table-layout: fixed; background-color: #DAD4CB; padding-bottom: 60px;">

		<table class="main" style="border-spacing: 0; width: 100%; max-width: 666px; background-color: #FDFBF6; font-family: 'Inter', Arial, Helvetica, sans-serif; color: #72513C; text-align: center;" width="100%">

		<tr class="logo-wrapper" style="background-color: #72513C; height: 66px;">
			<td style=""">
				<a href="https://www.ommindshop.com" target="_blank"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/logo.png" alt="Ommind" width="147" height="34" style=""></a>
			</td>
		</tr>
					
		<tr > 
			<td>
				<table align="center" role="presentation" width="100%" style="max-width: 500px;">
					<tr>
						<td align="center">
							<p style="font-weight: 600px; font-size: 30px;">Hi Beautiful Soul,</p>
							<p style="font-size: 14px;">Thank you again for taking our Chakra Test — we hope it brought you clarity on where you are in your inner journey.</p>
							<p style="font-size: 14px;">We’re now building something special:</p>
							<p style="font-size: 14px;">A personalized meditation and spiritual growth app inspired by Tibetan wisdom, designed to guide you step by step — from calm, to growth, to awakening.</p>
							<p style="font-size: 14px; margin-bottom: 0px;">We’d love your help.</p>
							<p style="font-size: 14px; margin-bottom: 0px; margin-top: 0px;"">Would you take <strong>2–3 minutes</strong> to answer a short survey that will shape the next version of <strong>OmMind?</strong></p>
							<p style="font-size: 14px; margin-bottom: 0px">As a thank-you:</p>
							<p style="font-size: 14px; margin-bottom: 0px; margin-top: 0px;">✨ 1 in 5 participants will receive a <strong>free chakra crystal bracelet</strong></p>
							<p style="font-size: 14px; margin-bottom: 0px; margin-top: 0px;">💎 Everyone will receive a <strong>50% discount</strong> for our OmMind Shop	</p>
							<p style="font-size: 14px;">Your voice truly matters — whether you're just starting out or have been walking this path for years.</p>
							<p style="font-size: 14px"><a target="_blank" href="https://docs.google.com/forms/d/e/1FAIpQLScy75LVfzwo_NJa-qAXDHaFL_vusNppajUl7Lg_VFy_6OphcQ/viewform">[Start the Survey]</a></p>
						</td>
					</tr>
					
				</table>
			</td>
		</tr>
		
				<tr > 
			<td>
				<table align="center" role="presentation" width="100%" style="max-width: 600px; margin-bottom: 20px;">
					<tr>
						<td align="center">
							<tr>
								<td height="234" width="601" align="center" valign="middle" background="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/xr.png" style="height: 234px; width: 601px; background-image: url('https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/xr.png'); background-repeat: no-repeat; background-position: center center; background-size: cover; text-align: center; vertical-align: middle;">
									<div>
										<table align="center" role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 500px;"> 
											<tr>
												<td align="center" style="padding: 20px 10px;"> 													
												</td>
											</tr>
										</table>
									</div>
								</td>
							</tr>
						</td>
					</tr>
					
				</table>
			</td>
		</tr>
		
		<tr>
			<td style="padding:0">
				<table width="100%" style="background-color: #72513C; border-spacing: 0px;">
					<tr>
						<td>
							<tr class="logo-wrapper">
								<td align="center" style="padding-top: 20px;">
									<a href="https://www.ommindshop.com/" target="_blank"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/logo.png" alt="Ommind" width="147" height="34" style=""></a>
								</td>
							</tr>						
						</td>
					</tr>
					<tr >
						<td align="center" style="padding: 25px 25px; font-size: 0;">
							<a href="https://www.facebook.com/profile.php?id=61573346513645" target="_blank" style="text-decoration: none; display: inline-block; vertical-align: middle;"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/company-info-logo-1.png"  alt="Social Media 1" width="25" height="25" style="width: 25px; height: 25px; border: 0; display: inline-block; margin: 0 12px; vertical-align: middle;"></a>
							<a href="https://www.tiktok.com/@ommindshop" target="_blank" style="text-decoration: none; display: inline-block; vertical-align: middle;"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/company-info-logo-2.png"  alt="Social Media 2" width="25" height="25" style="width: 25px; height: 25px; border: 0; display: inline-block; margin: 0 12px; vertical-align: middle;"></a>
							<a href="https://www.instagram.com/ommind_shop/" target="_blank" style="text-decoration: none; display: inline-block; vertical-align: middle;"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/company-info-logo-3.png"  alt="Social Media 3" width="25" height="25" style="width: 25px; height: 25px; border: 0; display: inline-block; margin: 0 12px; vertical-align: middle;"></a>
							<a href="https://www.youtube.com/@OmMind-Official" target="_blank" style="text-decoration: none; display: inline-block; vertical-align: middle;"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/company-info-logo-4.png"  alt="Social Media 4" width="25" height="25" style="width: 25px; height: 25px; border: 0; display: inline-block; margin: 0 12px; vertical-align: middle;"></a>
						</td>
					</tr>
					<tr>
						<td align="center" style="padding: 0 25px 50px 25px;"> 
							<p style="font-size: 14px; color: white; margin-bottom: 5px; margin-top: 0; margin-bottom: 10px;">86-90 Paul Street, London, EC2A 4UX, United Kingdom</p>
							<p style="font-size: 14px; color: #BDBDBD; margin-top: 0; margin-bottom: 0;">COPYRIGHT ©2025 OmMind Shop All Rights Reserved</p>
						</td>
					</tr>
					
					
				</table>
			</td>
		</tr>			
		</table> 

	</center> 
</body>
</html>
`
